package async

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"factura-movil-gateway/internal/domain"
)

type DTEDocument struct {
	XMLName  xml.Name    `xml:"DTE"`
	Version  string      `xml:"version,attr"`
	Document XMLDocument `xml:"Documento"`
}

type XMLDocument struct {
	XMLName   xml.Name    `xml:"Documento"`
	ID        string      `xml:"ID,attr"`
	Header    XMLHeader   `xml:"Encabezado"`
	Details   []XMLDetail `xml:"Detalle"`
	TED       XMLTED      `xml:"TED"`
	Timestamp string      `xml:"TmstFirma"`
}

type XMLHeader struct {
	XMLName  xml.Name    `xml:"Encabezado"`
	DocInfo  XMLDocInfo  `xml:"IdDoc"`
	Issuer   XMLIssuer   `xml:"Emisor"`
	Receiver XMLReceiver `xml:"Receptor"`
	Totals   XMLTotals   `xml:"Totales"`
}

type XMLDocInfo struct {
	XMLName      xml.Name `xml:"IdDoc"`
	DocumentType uint8    `xml:"TipoDTE"`
	Folio        int      `xml:"Folio"`
	IssueDate    string   `xml:"FchEmis"`
	PaymentForm  uint8    `xml:"FmaPago"`
	DueDate      string   `xml:"FchVenc"`
}

type XMLIssuer struct {
	XMLName      xml.Name `xml:"Emisor"`
	RUT          string   `xml:"RUTEmisor"`
	CompanyName  string   `xml:"RznSoc"`
	BusinessLine string   `xml:"GiroEmis"`
	Email        string   `xml:"CorreoEmisor"`
	Activities   []string `xml:"Acteco"`
	Address      string   `xml:"DirOrigen"`
	Commune      string   `xml:"CmnaOrigen"`
	City         string   `xml:"CiudadOrigen"`
}

type XMLReceiver struct {
	XMLName      xml.Name `xml:"Receptor"`
	RUT          string   `xml:"RUTRecep"`
	CompanyName  string   `xml:"RznSocRecep"`
	BusinessLine string   `xml:"GiroRecep"`
	Address      string   `xml:"DirRecep"`
	Commune      string   `xml:"CmnaRecep"`
	City         string   `xml:"CiudadRecep"`
}

type XMLTotals struct {
	XMLName   xml.Name `xml:"Totales"`
	NetAmount float64  `xml:"MntNeto"`
	TaxRate   float64  `xml:"TasaIVA"`
	TaxAmount float64  `xml:"IVA"`
	Total     float64  `xml:"MntTotal"`
}

type XMLDetail struct {
	XMLName     xml.Name    `xml:"Detalle"`
	LineNumber  int         `xml:"NroLinDet"`
	ItemCode    XMLItemCode `xml:"CdgItem"`
	ItemName    string      `xml:"NmbItem"`
	Description string      `xml:"DscItem"`
	Quantity    float64     `xml:"QtyItem"`
	Unit        string      `xml:"UnmdItem"`
	UnitPrice   float64     `xml:"PrcItem"`
	LineTotal   float64     `xml:"MontoItem"`
}

type XMLItemCode struct {
	XMLName   xml.Name `xml:"CdgItem"`
	CodeType  string   `xml:"TpoCodigo"`
	CodeValue string   `xml:"VlrCodigo"`
}

type XMLTED struct {
	XMLName xml.Name `xml:"TED"`
	Version string   `xml:"version,attr"`
	DD      XMLDD    `xml:"DD"`
	FRMT    XMLFRMT  `xml:"FRMT"`
}

type XMLDD struct {
	XMLName xml.Name        `xml:"DD"`
	RE      string          `xml:"RE"`
	TD      uint8           `xml:"TD"`
	F       int             `xml:"F"`
	FE      string          `xml:"FE"`
	RR      string          `xml:"RR"`
	RSR     string          `xml:"RSR"`
	MNT     uint64          `xml:"MNT"`
	IT1     string          `xml:"IT1"`
	CAF     domain.StampCAF `xml:"CAF"`
	TSTED   string          `xml:"TSTED"`
}

type XMLFRMT struct {
	XMLName   xml.Name `xml:"FRMT"`
	Algorithm string   `xml:"algoritmo,attr"`
	Value     string   `xml:",chardata"`
}

func ParseDTEXML(xmlData []byte) (*DTEDocument, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "iso-8859-1":
			return transform.NewReader(input, charmap.ISO8859_1.NewDecoder()), nil
		default:
			return nil, fmt.Errorf("unsupported charset: %s", charset)
		}
	}

	var dte DTEDocument
	err := decoder.Decode(&dte)
	if err != nil {
		return nil, err
	}
	return &dte, nil
}

func (dte *DTEDocument) ToInvoice() (*domain.Invoice, error) {
	doc := dte.Document

	issueDate, err := time.Parse("2006-01-02", doc.Header.DocInfo.IssueDate)
	if err != nil {
		issueDate = time.Now()
	}

	invoice := &domain.Invoice{
		DocumentType: doc.Header.DocInfo.DocumentType,
		Folio:        doc.Header.DocInfo.Folio,
		IssueDate:    issueDate,
		Issuer: domain.Company{
			Code:    doc.Header.Issuer.RUT,
			Name:    doc.Header.Issuer.CompanyName,
			Address: formatAddress(doc.Header.Issuer.Address, doc.Header.Issuer.Commune, doc.Header.Issuer.City),
		},
		Receiver: &domain.Company{
			Code:    doc.Header.Receiver.RUT,
			Name:    doc.Header.Receiver.CompanyName,
			Address: formatAddress(doc.Header.Receiver.Address, doc.Header.Receiver.Commune, doc.Header.Receiver.City),
		},
		Details: convertXMLDetails(doc.Details),
		Totals: domain.InvoiceTotals{
			TaxableAmount: doc.Header.Totals.NetAmount,
			TaxAmount:     doc.Header.Totals.TaxAmount,
			TotalAmount:   doc.Header.Totals.Total,
		},
	}

	return invoice, nil
}

func formatAddress(address, commune, city string) string {
	if address == "" {
		return ""
	}
	if commune != "" && city != "" {
		return address + ", " + commune + ", " + city
	}
	if commune != "" {
		return address + ", " + commune
	}
	if city != "" {
		return address + ", " + city
	}
	return address
}

func convertXMLDetails(xmlDetails []XMLDetail) []domain.InvoiceDetail {
	details := make([]domain.InvoiceDetail, len(xmlDetails))
	for i, detail := range xmlDetails {
		description := detail.ItemName
		if detail.Description != "" {
			description = detail.ItemName + " - " + detail.Description
		}

		details[i] = domain.InvoiceDetail{
			Quantity:    detail.Quantity,
			Description: description,
			UnitPrice:   detail.UnitPrice,
			LineTotal:   detail.LineTotal,
		}
	}
	return details
}
