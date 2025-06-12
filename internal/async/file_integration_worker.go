package async

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/usecases"
)

var _ Worker = &FileIntegrationWorker{}

type FileIntegrationWorker struct {
	ticker               *time.Ticker
	sourceDirectory      string
	inprogressDirectory  string
	destinationDirectory string
	errorDirectory       string
	documentService      usecases.DocumentService
}

type FileProcessingResult struct {
	OriginalFile   string
	StampFile      string
	PDF417File     string
	ThermalFile    string
	ProcessingTime time.Duration
	Error          error
}

// NewFileIntegrationWorker creates a new FileIntegrationWorker instance
func NewFileIntegrationWorker(
	tickerInterval time.Duration,
	sourceDirectory, inprogressDirectory, destinationDirectory, errorDirectory string,
	stampService usecases.StampService,
	companyService usecases.CompanyService,
) *FileIntegrationWorker {
	return &FileIntegrationWorker{
		ticker:               time.NewTicker(tickerInterval),
		sourceDirectory:      sourceDirectory,
		inprogressDirectory:  inprogressDirectory,
		destinationDirectory: destinationDirectory,
		errorDirectory:       errorDirectory,
		documentService:      usecases.NewDocumentService(stampService, companyService),
	}
}

func (w *FileIntegrationWorker) Run(ctx context.Context, done func()) {
	slog.Debug("file integration worker initialized",
		"sourceDir", w.sourceDirectory,
		"inprogressDir", w.inprogressDirectory,
		"destinationDir", w.destinationDirectory,
		"errorDir", w.errorDirectory)
	defer done()

	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			slog.Info("file integration worker cancelled, waiting for active processing to complete")
			wg.Wait()
			return
		case <-w.ticker.C:
			wg.Add(1)
			tickCtx := context.Background()
			go w.handleFileIntegration(tickCtx, wg.Done)
		}
	}
}

func (w *FileIntegrationWorker) handleFileIntegration(ctx context.Context, done func()) {
	defer done()

	slog.Debug("starting file integration tick")

	results, err := w.processAllDocuments()
	if err != nil {
		slog.Error("failed to process documents",
			"error", err,
			"sourceDir", w.sourceDirectory)
		return
	}

	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Error != nil {
			errorCount++
			slog.Error("document processing failed",
				"file", result.OriginalFile,
				"error", result.Error,
				"duration", result.ProcessingTime)
		} else {
			successCount++
			slog.Info("processing file", slog.String("filename", result.OriginalFile))
		}
	}

	if len(results) > 0 {
		slog.Info("batch complete",
			slog.Int("files_processed", len(results)),
			slog.Int("files_failed", errorCount),
		)
	} else {
		slog.Debug("file integration tick completed - no files to process")
	}
}

func (w *FileIntegrationWorker) processAllDocuments() ([]FileProcessingResult, error) {
	slog.Info("Starting document processing batch",
		"sourceDir", w.sourceDirectory,
		"inprogressDir", w.inprogressDirectory,
		"destinationDir", w.destinationDirectory,
		"errorDir", w.errorDirectory)

	if err := w.ensureDirectoriesExist(); err != nil {
		return nil, fmt.Errorf("failed to ensure directories exist: %w", err)
	}

	files, err := w.getSourceFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get source files: %w", err)
	}

	if len(files) == 0 {
		slog.Info("No files found in source directory")
		return []FileProcessingResult{}, nil
	}

	slog.Info("Found files to process", "count", len(files))

	var results []FileProcessingResult
	for _, file := range files {
		result := w.processDocument(file)
		results = append(results, result)

		if result.Error != nil {
			slog.Error("Failed to process document",
				"file", file,
				"error", result.Error)
		} else {
			slog.Info("Successfully processed document",
				"file", file,
				"processingTime", result.ProcessingTime)
		}
	}

	return results, nil
}

func (w *FileIntegrationWorker) processDocument(sourceFile string) FileProcessingResult {
	startTime := time.Now()

	result := FileProcessingResult{
		OriginalFile: sourceFile,
	}

	inProgressFile, err := w.moveToInProgress(sourceFile)
	if err != nil {
		result.Error = fmt.Errorf("failed to move file to in-progress: %w", err)
		return result
	}

	defer func() {
		if result.Error != nil {
			w.moveToError(inProgressFile, result.Error)
		}
	}()

	invoice, err := w.parseXMLToInvoice(inProgressFile)
	if err != nil {
		result.Error = fmt.Errorf("failed to parse XML to invoice: %w", err)
		return result
	}

	processingResult, err := w.documentService.ProcessInvoice(invoice)
	if err != nil {
		result.Error = fmt.Errorf("failed to process invoice: %w", err)
		return result
	}

	err = w.saveFilesToDestination(inProgressFile, processingResult, &result)
	if err != nil {
		result.Error = fmt.Errorf("failed to save files to destination: %w", err)
		return result
	}

	result.ProcessingTime = time.Since(startTime)
	return result
}

func (w *FileIntegrationWorker) moveToInProgress(sourceFile string) (string, error) {
	fileName := filepath.Base(sourceFile)
	inProgressFile := filepath.Join(w.inprogressDirectory, fileName)

	if err := w.moveFile(sourceFile, inProgressFile); err != nil {
		return "", err
	}

	slog.Debug("Moved file to in-progress",
		"from", sourceFile,
		"to", inProgressFile)

	return inProgressFile, nil
}

func (w *FileIntegrationWorker) moveToError(inProgressFile string, processingError error) {
	fileName := filepath.Base(inProgressFile)
	errorFile := filepath.Join(w.errorDirectory, fileName)

	if err := w.moveFile(inProgressFile, errorFile); err != nil {
		slog.Error("Failed to move file to error directory",
			"file", inProgressFile,
			"errorDir", w.errorDirectory,
			"moveError", err,
			"originalError", processingError)
		os.Remove(inProgressFile)
		return
	}

	slog.Info("Moved failed file to error directory",
		"from", inProgressFile,
		"to", errorFile,
		"error", processingError)
}

func (w *FileIntegrationWorker) parseXMLToInvoice(filePath string) (*domain.Invoice, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	dte, err := ParseDTEXML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DTE XML: %w", err)
	}

	invoice, err := dte.ToInvoice()
	if err != nil {
		return nil, fmt.Errorf("failed to convert DTE to invoice: %w", err)
	}

	slog.Debug("Parsed XML to invoice",
		"documentType", invoice.DocumentType,
		"folio", invoice.Folio,
		"issuer", invoice.Issuer.Name)

	return invoice, nil
}

func (w *FileIntegrationWorker) saveFilesToDestination(inProgressFile string, processingResult usecases.ProcessingResult, result *FileProcessingResult) error {
	baseName := strings.TrimSuffix(filepath.Base(inProgressFile), filepath.Ext(inProgressFile))

	originalDest := filepath.Join(w.destinationDirectory, filepath.Base(inProgressFile))
	if err := w.moveFile(inProgressFile, originalDest); err != nil {
		return fmt.Errorf("failed to move original file: %w", err)
	}
	result.OriginalFile = originalDest

	stampFile := filepath.Join(w.destinationDirectory, baseName+"_stamp.xml")
	if err := os.WriteFile(stampFile, processingResult.StampXML, 0644); err != nil {
		return fmt.Errorf("failed to save stamp file: %w", err)
	}
	result.StampFile = stampFile

	pdf417File := filepath.Join(w.destinationDirectory, baseName+"_pdf417.png")
	if err := os.WriteFile(pdf417File, processingResult.PDF417Data, 0644); err != nil {
		return fmt.Errorf("failed to save PDF417 file: %w", err)
	}
	result.PDF417File = pdf417File

	thermalFile := filepath.Join(w.destinationDirectory, baseName+"_thermal.pdf")
	if err := os.WriteFile(thermalFile, processingResult.ThermalPDF, 0644); err != nil {
		return fmt.Errorf("failed to save thermal PDF file: %w", err)
	}
	result.ThermalFile = thermalFile

	slog.Debug("Saved all files to destination",
		"original", result.OriginalFile,
		"stamp", result.StampFile,
		"pdf417", result.PDF417File,
		"thermal", result.ThermalFile)

	return nil
}

func (w *FileIntegrationWorker) moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Remove(src)
}

func (w *FileIntegrationWorker) ensureDirectoriesExist() error {
	dirs := []string{w.sourceDirectory, w.inprogressDirectory, w.destinationDirectory, w.errorDirectory}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (w *FileIntegrationWorker) getSourceFiles() ([]string, error) {
	entries, err := os.ReadDir(w.sourceDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if strings.HasSuffix(strings.ToLower(fileName), ".xml") {
			fullPath := filepath.Join(w.sourceDirectory, fileName)
			files = append(files, fullPath)
		}
	}

	return files, nil
}

func (w *FileIntegrationWorker) Shutdown() {
	slog.Info("shutting down file integration worker")
	w.ticker.Stop()
}
