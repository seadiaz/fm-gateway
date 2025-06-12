package usecases

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"strings"
	"time"

	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/utils"

	"github.com/jung-kurt/gofpdf/v2"
)

// DocumentService defines the interface for document processing operations
type DocumentService interface {
	ProcessInvoice(invoice *domain.Invoice) (ProcessingResult, error)
}

// ProcessingResult contains the results of document processing
type ProcessingResult struct {
	StampXML       []byte
	PDF417Data     []byte
	ThermalPDF     []byte
	ProcessingTime time.Duration
	Error          error
}

// SimpleDocumentService implements DocumentService
type SimpleDocumentService struct {
	stampService   StampService
	companyService CompanyService
}

// NewDocumentService creates a new document service
func NewDocumentService(stampService StampService, companyService CompanyService) DocumentService {
	return &SimpleDocumentService{
		stampService:   stampService,
		companyService: companyService,
	}
}

// ProcessInvoice processes a single document through the complete workflow
func (s *SimpleDocumentService) ProcessInvoice(invoice *domain.Invoice) (ProcessingResult, error) {
	startTime := time.Now()

	result := ProcessingResult{}

	stampXML, err := s.createStamp(invoice)
	if err != nil {
		result.Error = fmt.Errorf("failed to create stamp: %w", err)
		return result, result.Error
	}
	result.StampXML = stampXML

	pdf417Data, err := s.createPDF417(stampXML)
	if err != nil {
		result.Error = fmt.Errorf("failed to create PDF417: %w", err)
		return result, result.Error
	}
	result.PDF417Data = pdf417Data

	thermalPDF, err := s.createThermalPDF(invoice, stampXML)
	if err != nil {
		result.Error = fmt.Errorf("failed to create thermal PDF: %w", err)
		return result, result.Error
	}
	result.ThermalPDF = thermalPDF

	result.ProcessingTime = time.Since(startTime)
	return result, nil
}

// createStamp creates a stamp for the invoice using StampService
func (s *SimpleDocumentService) createStamp(invoice *domain.Invoice) ([]byte, error) {
	ctx := context.Background()

	company, err := s.companyService.FindByCode(ctx, invoice.Issuer.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to find company with code %s: %w", invoice.Issuer.Code, err)
	}

	stamp, err := s.stampService.Generate(ctx, *company, *invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to generate stamp: %w", err)
	}

	stampXML, err := s.convertStampToXML(stamp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert stamp to XML: %w", err)
	}

	slog.Debug("Created stamp using StampService", "size", len(stampXML), "company", company.Name)
	return stampXML, nil
}

// convertToCompany converts Invoice issuer to domain.Company
func (s *SimpleDocumentService) convertToCompany(invoice *domain.Invoice) domain.Company {
	return domain.Company{
		ID:      "temp-company-id",
		Code:    invoice.Issuer.Code,
		Name:    invoice.Issuer.Name,
		Address: invoice.Issuer.Address,
	}
}

// convertStampToXML converts domain.Stamp to XML bytes
func (s *SimpleDocumentService) convertStampToXML(stamp domain.Stamp) ([]byte, error) {
	// Create TED structure for XML serialization
	ted := struct {
		XMLName xml.Name  `xml:"TED"`
		Version string    `xml:"version,attr"`
		DD      domain.DD `xml:"DD"`
		FRMT    struct {
			Algorithm string `xml:"algoritmo,attr"`
			Value     string `xml:",chardata"`
		} `xml:"FRMT"`
	}{
		Version: "1.0",
		DD:      stamp.DD,
		FRMT: struct {
			Algorithm string `xml:"algoritmo,attr"`
			Value     string `xml:",chardata"`
		}{
			Algorithm: "SHA1withRSA",
			Value:     stamp.FRMT,
		},
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(ted, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TED to XML: %w", err)
	}

	// Add XML header
	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fullXML := append(xmlHeader, xmlData...)

	return fullXML, nil
}

// createPDF417 creates a PDF417 barcode image from stamp data
func (s *SimpleDocumentService) createPDF417(stampXML []byte) ([]byte, error) {
	// Generate PDF417 with auto-dimensions
	_, pdf417Data, err := utils.GenerateStampPDF417FromXMLWithAutoDimensions(string(stampXML))
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF417: %w", err)
	}

	slog.Debug("Created PDF417", "size", len(pdf417Data))
	return pdf417Data, nil
}

// createThermalPDF creates a thermal printer optimized PDF document
func (s *SimpleDocumentService) createThermalPDF(invoice *domain.Invoice, stampXML []byte) ([]byte, error) {
	// Create PDF with custom page size for 3-inch thermal printer (80mm width)
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: 80,  // 3 inch width
			Ht: 200, // Initial height, will extend as needed
		},
		FontDirStr:     "",
		OrientationStr: "P",
	})

	pdf.AddPage()

	// Set smaller margins for thermal printer
	pdf.SetMargins(3, 3, 3)
	pdf.SetAutoPageBreak(true, 3)

	// Calculate separator width for 80mm thermal printer (approximately 48 characters)
	separatorWidth := 48

	// Helper function to format currency with thousand separators
	formatCurrency := func(amount float64) string {
		// Format with thousand separators using dots (Chilean format)
		formatted := fmt.Sprintf("%.0f", amount)

		// Add thousand separators (dots)
		if len(formatted) > 3 {
			var result []rune
			for i, digit := range []rune(formatted) {
				if i > 0 && (len(formatted)-i)%3 == 0 {
					result = append(result, '.')
				}
				result = append(result, digit)
			}
			formatted = string(result)
		}

		return "$" + formatted
	}

	// Helper function to convert Spanish characters to ASCII equivalents for PDF compatibility
	encodeText := func(text string) string {
		// Replace Spanish accented characters with ASCII equivalents
		replacements := map[string]string{
			"á": "a", "à": "a", "ä": "a", "â": "a",
			"é": "e", "è": "e", "ë": "e", "ê": "e",
			"í": "i", "ì": "i", "ï": "i", "î": "i",
			"ó": "o", "ò": "o", "ö": "o", "ô": "o",
			"ú": "u", "ù": "u", "ü": "u", "û": "u",
			"ñ": "n",
			"Á": "A", "À": "A", "Ä": "A", "Â": "A",
			"É": "E", "È": "E", "Ë": "E", "Ê": "E",
			"Í": "I", "Ì": "I", "Ï": "I", "Î": "I",
			"Ó": "O", "Ò": "O", "Ö": "O", "Ô": "O",
			"Ú": "U", "Ù": "U", "Ü": "U", "Û": "U",
			"Ñ": "N",
			// Special characters
			"°": "o", // Degree symbol to lowercase 'o'
			"º": "o", // Masculine ordinal indicator
			"ª": "a", // Feminine ordinal indicator
		}

		result := text
		for accented, ascii := range replacements {
			result = strings.ReplaceAll(result, accented, ascii)
		}
		return result
	}

	// Company Header Section
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, encodeText(invoice.Issuer.Name), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, fmt.Sprintf("RUT: %s", invoice.Issuer.Code), "", 1, "C", false, 0, "")

	if invoice.Issuer.Address != "" {
		// Split long addresses
		address := encodeText(invoice.Issuer.Address)
		if len(address) > 35 {
			words := strings.Fields(address)
			var line1, line2 string
			for i, word := range words {
				if len(line1+" "+word) <= 35 {
					if line1 == "" {
						line1 = word
					} else {
						line1 += " " + word
					}
				} else {
					line2 = strings.Join(words[i:], " ")
					break
				}
			}
			pdf.CellFormat(0, 4, line1, "", 1, "C", false, 0, "")
			if line2 != "" {
				pdf.CellFormat(0, 4, line2, "", 1, "C", false, 0, "")
			}
		} else {
			pdf.CellFormat(0, 4, address, "", 1, "C", false, 0, "")
		}
	}

	pdf.Ln(2)

	// Document Type and Number
	pdf.SetFont("Arial", "B", 9)
	docTypeName := getDocumentTypeName(invoice.DocumentType)
	pdf.CellFormat(0, 5, encodeText(docTypeName), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("No. %d", invoice.Folio), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, fmt.Sprintf("Fecha: %s", invoice.IssueDate.Format("02/01/2006")), "", 1, "C", false, 0, "")

	// Separator line
	pdf.Ln(2)
	pdf.CellFormat(0, 3, strings.Repeat("-", separatorWidth), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Customer Information
	if invoice.Receiver != nil {
		pdf.SetFont("Arial", "B", 8)
		pdf.CellFormat(0, 4, "CLIENTE:", "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 7)
		// Customer name with wrapping
		customerName := encodeText(invoice.Receiver.Name)
		if len(customerName) > 40 {
			words := strings.Fields(customerName)
			var line1, line2 string
			for i, word := range words {
				if len(line1+" "+word) <= 40 {
					if line1 == "" {
						line1 = word
					} else {
						line1 += " " + word
					}
				} else {
					line2 = strings.Join(words[i:], " ")
					break
				}
			}
			pdf.CellFormat(0, 3, line1, "", 1, "L", false, 0, "")
			if line2 != "" {
				pdf.CellFormat(0, 3, line2, "", 1, "L", false, 0, "")
			}
		} else {
			pdf.CellFormat(0, 3, customerName, "", 1, "L", false, 0, "")
		}

		if invoice.Receiver.Code != "" {
			pdf.CellFormat(0, 3, fmt.Sprintf("RUT: %s", invoice.Receiver.Code), "", 1, "L", false, 0, "")
		}

		pdf.Ln(2)
		pdf.CellFormat(0, 3, strings.Repeat("-", separatorWidth), "", 1, "C", false, 0, "")
		pdf.Ln(2)
	}

	// Items Header
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(0, 3, "DETALLE", "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Items
	pdf.SetFont("Arial", "", 7)
	for _, item := range invoice.Details {
		// Item description (truncate if too long)
		desc := encodeText(item.Description)
		if len(desc) > 38 {
			desc = desc[:35] + "..."
		}
		pdf.CellFormat(0, 3, desc, "", 1, "L", false, 0, "")

		// Quantity, unit price, and total on separate line
		pdf.CellFormat(0, 3, fmt.Sprintf("%.0f x %s", item.Quantity, formatCurrency(item.UnitPrice)), "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 3, fmt.Sprintf("%s", formatCurrency(item.LineTotal)), "", 1, "R", false, 0, "")
		pdf.Ln(1)
	}

	// Totals separator
	pdf.Ln(1)
	pdf.CellFormat(0, 3, strings.Repeat("-", separatorWidth), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Totals Section
	pdf.SetFont("Arial", "", 8)

	// Add separator before totals
	pdf.CellFormat(0, 3, strings.Repeat("-", separatorWidth), "", 1, "C", false, 0, "")
	pdf.Ln(1)

	if invoice.Totals.TaxableAmount > 0 {
		pdf.CellFormat(0, 4, "Subtotal:", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 4, fmt.Sprintf("%s", formatCurrency(invoice.Totals.TaxableAmount)), "", 1, "R", false, 0, "")
	}

	if invoice.Totals.TaxAmount > 0 {
		pdf.CellFormat(0, 4, "IVA (19%):", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 4, fmt.Sprintf("%s", formatCurrency(invoice.Totals.TaxAmount)), "", 1, "R", false, 0, "")
	}

	// Total amount (highlighted)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 5, "TOTAL:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("%s", formatCurrency(invoice.Totals.TotalAmount)), "", 1, "R", false, 0, "")

	// Final separator
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 3, strings.Repeat("=", separatorWidth), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// SII Electronic Stamp Section
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 4, encodeText("TIMBRE ELECTRONICO SII"), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 3, "Res. 80 de 2014", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 3, "Verifique en www.sii.cl", "", 1, "C", false, 0, "")

	pdf.Ln(3)

	// Generate and embed PDF417 barcode
	pdf417Image, _, err := utils.GenerateStampPDF417FromXMLWithAutoDimensions(string(stampXML))
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF417 for embedding: %w", err)
	}

	// Convert to 8-bit PNG for gofpdf compatibility
	pdf417Data, err := convertTo8BitPNG(pdf417Image)
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF417 to 8-bit PNG: %w", err)
	}

	// Register the PDF417 image from memory
	imageInfo := pdf.RegisterImageOptionsReader("pdf417", gofpdf.ImageOptions{
		ImageType: "PNG",
	}, bytes.NewReader(pdf417Data))

	if pdf.Error() != nil {
		return nil, fmt.Errorf("failed to register PDF417 image: %w", pdf.Error())
	}

	// Calculate image dimensions for thermal printer - make it prominent
	pageWidth, _ := pdf.GetPageSize()
	margins := 6.0 // Margins for barcode
	maxWidth := pageWidth - margins

	// Make the barcode large and prominent
	imgWidth := maxWidth * 0.9 // Use 90% of available width

	// Calculate height maintaining aspect ratio
	originalWidth := float64(imageInfo.Width())
	originalHeight := float64(imageInfo.Height())
	aspectRatio := originalHeight / originalWidth
	imgHeight := imgWidth * aspectRatio

	// Center the image horizontally
	x := (pageWidth - imgWidth) / 2

	// Add the PDF417 barcode image
	pdf.ImageOptions("pdf417", x, pdf.GetY(), imgWidth, imgHeight, false, gofpdf.ImageOptions{}, 0, "")

	// Move cursor below the image
	pdf.SetY(pdf.GetY() + imgHeight + 3)

	// Convert PDF to bytes
	var buf strings.Builder
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	pdfBytes := []byte(buf.String())
	slog.Debug("Created thermal receipt PDF with embedded PDF417", "size", len(pdfBytes))
	return pdfBytes, nil
}

// getDocumentTypeName returns the human-readable name for a document type
func getDocumentTypeName(docType uint8) string {
	switch docType {
	case 33:
		return "FACTURA ELECTRONICA"
	case 34:
		return "FACTURA NO AFECTA O EXENTA"
	case 39:
		return "BOLETA ELECTRONICA"
	case 41:
		return "BOLETA EXENTA ELECTRONICA"
	case 43:
		return "LIQUIDACION FACTURA"
	case 46:
		return "FACTURA DE COMPRA"
	case 52:
		return "GUIA DE DESPACHO"
	case 56:
		return "NOTA DE DEBITO"
	case 61:
		return "NOTA DE CREDITO"
	case 110:
		return "FACTURA DE EXPORTACION"
	case 111:
		return "NOTA DE DEBITO EXPORTACION"
	case 112:
		return "NOTA DE CREDITO EXPORTACION"
	default:
		return fmt.Sprintf("DOCUMENTO TIPO %d", docType)
	}
}

// convertTo8BitPNG converts an image to 8-bit PNG format for gofpdf compatibility
func convertTo8BitPNG(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, color.Palette{
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{0, 0, 0, 255},       // Black
	})

	// Convert to paletted image (8-bit)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()

			// Convert to grayscale and then to black/white
			gray := (r + g + b) / 3
			if gray > 32768 { // Threshold for white vs black
				paletted.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				paletted.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, paletted)
	if err != nil {
		return nil, fmt.Errorf("failed to encode 8-bit PNG: %w", err)
	}

	return buf.Bytes(), nil
}
