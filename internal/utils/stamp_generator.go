package utils

import (
	"encoding/xml"
	"fmt"

	"factura-movil-gateway/internal/domain"
)

// TED represents the main Timbre Electrónico structure
type TED struct {
	XMLName xml.Name `xml:"TED"`
	Version string   `xml:"version,attr"`
	DD      DD       `xml:"DD"`
	FRMT    FRMT     `xml:"FRMT"`
}

// DD represents the document data
type DD struct {
	XMLName xml.Name `xml:"DD"`
	RE      string   `xml:"RE"`    // RUT Emisor
	TD      uint8    `xml:"TD"`    // Tipo Documento
	F       int      `xml:"F"`     // Folio
	FE      string   `xml:"FE"`    // Fecha Emisión
	RR      string   `xml:"RR"`    // RUT Receptor
	RSR     string   `xml:"RSR"`   // Razón Social Receptor
	MNT     int      `xml:"MNT"`   // Monto Total
	IT1     string   `xml:"IT1"`   // Item 1
	CAF     CAFRef   `xml:"CAF"`   // CAF Reference
	TSTED   string   `xml:"TSTED"` // Timestamp
}

// CAFRef represents a reference to the CAF used
type CAFRef struct {
	XMLName xml.Name `xml:"CAF"`
	Version string   `xml:"version,attr"`
	DA      DA       `xml:"DA"`
	FRMA    string   `xml:"FRMA"`
}

// DA represents document authorization data
type DA struct {
	XMLName xml.Name `xml:"DA"`
	RE      string   `xml:"RE"`  // RUT Emisor
	RS      string   `xml:"RS"`  // Razón Social
	TD      uint8    `xml:"TD"`  // Tipo Documento
	RNG     RNG      `xml:"RNG"` // Rango
	FA      string   `xml:"FA"`  // Fecha Autorización
	RSAPK   RSAPK    `xml:"RSAPK"`
	IDK     int      `xml:"IDK"`
}

// RNG represents the folio range
type RNG struct {
	XMLName xml.Name `xml:"RNG"`
	D       int      `xml:"D"` // Desde
	H       int      `xml:"H"` // Hasta
}

// RSAPK represents RSA public key
type RSAPK struct {
	XMLName xml.Name `xml:"RSAPK"`
	M       string   `xml:"M"` // Modulus
	E       string   `xml:"E"` // Exponent
}

// FRMT represents the signature format
type FRMT struct {
	XMLName   xml.Name `xml:"FRMT"`
	Algorithm string   `xml:"algoritmo,attr"`
	Value     string   `xml:",chardata"`
}

// GenerateStampXML generates a stamp XML from stamp data
func GenerateStampXML(stampData *domain.StampData) ([]byte, error) {
	// Create TED structure
	ted := TED{
		Version: "1.0",
		DD: DD{
			RE:  stampData.RutEmisor,
			TD:  stampData.TipoDoc,
			F:   stampData.Folio,
			FE:  stampData.FechaEmision,
			RR:  stampData.RutReceptor,
			RSR: "Cliente Ejemplo", // TODO: Get from stamp data
			MNT: stampData.MontoTotal,
			IT1: "Producto Ejemplo", // TODO: Get from stamp data
			CAF: CAFRef{
				Version: "1.0",
				DA: DA{
					RE: stampData.RutEmisor,
					RS: "Empresa Ejemplo S.A.", // TODO: Get from stamp data
					TD: stampData.TipoDoc,
					RNG: RNG{
						D: 1,    // TODO: Get from CAF
						H: 1000, // TODO: Get from CAF
					},
					FA: "2024-01-01", // TODO: Get from CAF
					RSAPK: RSAPK{
						M: "sampleModulus", // TODO: Get from CAF
						E: "AQAB",          // Standard RSA exponent
					},
					IDK: 1, // TODO: Get from CAF
				},
				FRMA: "sampleSignature", // TODO: Generate actual signature
			},
			TSTED: "2024-01-01T00:00:00", // TODO: Generate timestamp
		},
		FRMT: FRMT{
			Algorithm: "SHA1withRSA",
			Value:     "sampleStampSignature", // TODO: Generate actual signature
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
