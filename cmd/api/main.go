package main

import (
	"context"
	"factura-movil-gateway/internal/async"
	"factura-movil-gateway/internal/controllers"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/persistence"
	"factura-movil-gateway/internal/storage"
	"factura-movil-gateway/internal/usecases"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug, ReplaceAttr: slogReplaceSource})))
	slog.Info("🚀 fm gateway is initializing")

	dbhost := os.Getenv("FMG_DBHOST")
	dbuser := os.Getenv("FMG_DBUSER")
	dbpass := os.Getenv("FMG_DBPASS")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable", dbhost, dbuser, dbpass)

	storage := storage.NewLocalStorage("tmp")
	cafRepository, err := persistence.NewCAFRepository(dsn)
	if err != nil {
		panic(err)
	}
	cafService := usecases.NewCAFService(storage, cafRepository)

	companyRepository, err := persistence.NewCompanyRepository(dsn)
	if err != nil {
		panic(err)
	}
	companyService := usecases.NewCompanyService(companyRepository)

	stampService := usecases.NewStampService(cafService)

	httpServer := httpserver.NewServer(
		controllers.NewCAFController(cafService, companyService),
		controllers.NewStampController(stampService, companyService),
		controllers.NewCompanyController(companyService),
	)

	ctx, cancelFn := context.WithCancel(context.Background())

	go httpServer.Run()

	sourceDir := getEnvOrDefault("FMG_PROCESSOR_SOURCE_DIR", "./tmp/source")
	inprogressDir := getEnvOrDefault("FMG_PROCESSOR_INPROGRESS_DIR", "./tmp/inprogress")
	destinationDir := getEnvOrDefault("FMG_PROCESSOR_DESTINATION_DIR", "./tmp/destination")
	errorDir := getEnvOrDefault("FMG_PROCESSOR_ERROR_DIR", "./tmp/errors")
	intervalStr := getEnvOrDefault("FMG_PROCESSOR_INTERVAL", "5s")

	var fileWorker *async.FileIntegrationWorker
	if sourceDir == "" || inprogressDir == "" || destinationDir == "" || errorDir == "" {
		panic("FMG_PROCESSOR_SOURCE_DIR, FMG_PROCESSOR_INPROGRESS_DIR, FMG_PROCESSOR_DESTINATION_DIR, and FMG_PROCESSOR_ERROR_DIR must be set")
	}
	// Parse interval
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		slog.Warn("Invalid processor interval, using default 30s",
			"provided", intervalStr,
			"error", err)
		interval = 30 * time.Second
	}

	// Ensure directories exist
	if err := ensureDirectoriesExist(sourceDir, inprogressDir, destinationDir, errorDir); err != nil {
		slog.Error("Failed to create processor directories", "error", err)
		os.Exit(1)
	}

	// Create file integration worker
	fileWorker = async.NewFileIntegrationWorker(
		interval,
		sourceDir,
		inprogressDir,
		destinationDir,
		errorDir,
		stampService,
		companyService,
	)

	var wg sync.WaitGroup
	go fileWorker.Run(ctx, wg.Done)

	slog.Info("✅ File integration worker started successfully")

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	slog.Info("🌟 FM Gateway started successfully!")

	<-signalChannel
	slog.Info("👋 Shutdown signal received, stopping services...")

	// Cancel context to stop workers
	cancelFn()

	wg.Wait()

	slog.Info("✅ All services stopped. Good bye!")
	os.Exit(0)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ensureDirectoriesExist(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func slogReplaceSource(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}
