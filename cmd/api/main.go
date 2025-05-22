package main

import (
	"context"
	"factura-movil-gateway/internal/controllers"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/persistence"
	"factura-movil-gateway/internal/storage"
	"factura-movil-gateway/internal/usecases"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug, ReplaceAttr: slogReplaceSource})))
	slog.Info("🚀 zensor is initializing")

	storage := storage.NewLocalStorage("tmp")
	cafRepository := persistence.NewCAFRepository(storage)
	cafService := usecases.NewCAFService(cafRepository)

	httpServer := httpserver.NewServer(
		controllers.NewCAFController(cafService),
	)

	_, cancelFn := context.WithCancel(context.Background())
	go httpServer.Run()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel
	cancelFn()
	slog.Info("good bye!!!")
	os.Exit(0)
}

func slogReplaceSource(groups []string, a slog.Attr) slog.Attr {
	// Check if the attribute is the source key
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		// Set the file attribute to only its base name
		source.File = filepath.Base(source.File)
		return slog.Any(a.Key, source)
	}
	return a // Return unchanged attribute for others
}
