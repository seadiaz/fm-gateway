package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/pdf417"
)

// StampPDF417Data represents the structure that will be encoded in the PDF417 barcode
// This follows the Chilean electronic invoicing standard (TED structure)
type StampPDF417Data struct {
	XMLName xml.Name   `xml:"TED"`
	Version string     `xml:"version,attr"`
	DD      PDF417DD   `xml:"DD"`
	FRMT    PDF417FRMT `xml:"FRMT"`
}

type PDF417DD struct {
	RE    string    `xml:"RE"`    // Company RUT
	TD    uint8     `xml:"TD"`    // Document Type
	F     int64     `xml:"F"`     // Folio
	FE    string    `xml:"FE"`    // Document Date
	RR    string    `xml:"RR"`    // Customer RUT
	RSR   string    `xml:"RSR"`   // Customer Name
	MNT   uint64    `xml:"MNT"`   // Total Amount
	IT1   string    `xml:"IT1"`   // First Item Description
	CAF   PDF417CAF `xml:"CAF"`   // CAF Reference
	TSTED string    `xml:"TSTED"` // Timestamp
}

type PDF417CAF struct {
	Version string     `xml:"version,attr"`
	DA      PDF417DA   `xml:"DA"`
	FRMA    PDF417FRMA `xml:"FRMA"`
}

type PDF417DA struct {
	RE    string      `xml:"RE"`
	RS    string      `xml:"RS"`
	TD    uint8       `xml:"TD"`
	RNG   PDF417RNG   `xml:"RNG"`
	FA    string      `xml:"FA"`
	RSAPK PDF417RSAPK `xml:"RSAPK"`
	IDK   string      `xml:"IDK"`
}

type PDF417RNG struct {
	D int64 `xml:"D"`
	H int64 `xml:"H"`
}

type PDF417RSAPK struct {
	M string `xml:"M"`
	E string `xml:"E"`
}

type PDF417FRMA struct {
	Algorithm string `xml:"algoritmo,attr"`
	Value     string `xml:",chardata"`
}

type PDF417FRMT struct {
	Algorithm string `xml:"algoritmo,attr"`
	Value     string `xml:",chardata"`
}

// GenerateStampPDF417 creates a PDF417 barcode from stamp data
// Returns the barcode as an image and as raw bytes
func GenerateStampPDF417(stampData StampPDF417Data, width, height int) (image.Image, []byte, error) {
	// Serialize the stamp data to XML
	xmlData, err := xml.Marshal(stampData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal stamp data to XML: %w", err)
	}

	// Generate PDF417 barcode with security level 2 (recommended for documents)
	qrCode, err := pdf417.Encode(string(xmlData), 2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode PDF417: %w", err)
	}

	// Scale the barcode to the desired size
	scaledBarcode, err := barcode.Scale(qrCode, width, height)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scale PDF417: %w", err)
	}

	// Convert to PNG bytes
	var buf bytes.Buffer
	err = png.Encode(&buf, scaledBarcode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return scaledBarcode, buf.Bytes(), nil
}

// GenerateStampPDF417FromXML creates a PDF417 barcode directly from XML string
// This is useful when you already have the TED XML data
func GenerateStampPDF417FromXML(xmlData string, width, height int) (image.Image, []byte, error) {
	// Generate PDF417 barcode directly from XML string with security level 2
	qrCode, err := pdf417.Encode(xmlData, 2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode PDF417: %w", err)
	}

	// Scale the barcode to the desired size
	scaledBarcode, err := barcode.Scale(qrCode, width, height)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scale PDF417: %w", err)
	}

	// Convert to PNG bytes
	var buf bytes.Buffer
	err = png.Encode(&buf, scaledBarcode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return scaledBarcode, buf.Bytes(), nil
}

// ConvertDomainStampToPDF417Data converts a domain stamp to PDF417 data structure
func ConvertDomainStampToPDF417Data(re string, td uint8, f int64, fe string, rr string, rsr string, mnt uint64, it1 string, caf PDF417CAF, tsted string, frmt string) StampPDF417Data {
	return StampPDF417Data{
		Version: "1.0",
		DD: PDF417DD{
			RE:    re,
			TD:    td,
			F:     f,
			FE:    fe,
			RR:    rr,
			RSR:   rsr,
			MNT:   mnt,
			IT1:   it1,
			CAF:   caf,
			TSTED: tsted,
		},
		FRMT: PDF417FRMT{
			Algorithm: "SHA1withRSA",
			Value:     frmt,
		},
	}
}

// CalculateSafePDF417Dimensions determines appropriate dimensions for PDF417 based on data size
// This helps prevent "can not scale barcode to an image smaller than" errors
func CalculateSafePDF417Dimensions(dataLength int) (width, height int) {
	// Conservative base dimensions that work for most stamp data
	baseWidth := 500
	baseHeight := 200

	// Adjust dimensions based on data length with generous margins
	if dataLength > 2000 {
		// Large data needs much more space
		baseWidth = 700
		baseHeight = 300
	} else if dataLength > 1500 {
		// Medium-large data
		baseWidth = 600
		baseHeight = 250
	} else if dataLength > 1000 {
		// Standard data
		baseWidth = 550
		baseHeight = 220
	} else if dataLength > 500 {
		// Small-medium data
		baseWidth = 500
		baseHeight = 200
	} else {
		// Small data but still safe margins
		baseWidth = 450
		baseHeight = 180
	}

	return baseWidth, baseHeight
}

// GenerateStampPDF417WithAutoDimensions creates a PDF417 barcode with automatically calculated dimensions
// This is the recommended function to use as it prevents scaling errors
func GenerateStampPDF417WithAutoDimensions(stampData StampPDF417Data) (image.Image, []byte, error) {
	// Serialize the stamp data to XML to calculate size
	xmlData, err := xml.Marshal(stampData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal stamp data to XML: %w", err)
	}

	// Calculate safe dimensions
	width, height := CalculateSafePDF417Dimensions(len(xmlData))

	// Generate the barcode
	return GenerateStampPDF417(stampData, width, height)
}

// GenerateStampPDF417FromXMLWithAutoDimensions creates a PDF417 barcode with automatically calculated dimensions
// This is the recommended function to use as it prevents scaling errors
func GenerateStampPDF417FromXMLWithAutoDimensions(xmlData string) (image.Image, []byte, error) {
	// Calculate safe dimensions based on XML data length
	width, height := CalculateSafePDF417Dimensions(len(xmlData))

	// Generate the barcode
	return GenerateStampPDF417FromXML(xmlData, width, height)
}
