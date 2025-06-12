package async

import (
	"strings"
	"testing"
	"time"
)

func TestParseDTEXML(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="ISO-8859-1" standalone="no"?>
<DTE version="1.0">
  <Documento ID="DOC_29_33_2404">
    <Encabezado>
      <IdDoc>
        <TipoDTE>33</TipoDTE>
        <Folio>2404</Folio>
        <FchEmis>2025-05-05</FchEmis>
        <FmaPago>2</FmaPago>
        <FchVenc>2025-05-31</FchVenc>
      </IdDoc>
      <Emisor>
        <RUTEmisor>76212889-6</RUTEmisor>
        <RznSoc>FACTURA MOVIL SPA</RznSoc>
        <GiroEmis>COMERCIO AL POR MENOR DE COMPUTADORAS, SOFTWARES Y SUMINISTROS</GiroEmis>
        <CorreoEmisor>rfernandez@facturamovil.cl</CorreoEmisor>
        <Acteco>523930</Acteco>
        <Acteco>726000</Acteco>
        <DirOrigen>Vicuña Mackenna 9705</DirOrigen>
        <CmnaOrigen>La Florida</CmnaOrigen>
        <CiudadOrigen>Santiago</CiudadOrigen>
      </Emisor>
      <Receptor>
        <RUTRecep>77371419-3</RUTRecep>
        <RznSocRecep>AGRICOLA PAINE LTDA</RznSocRecep>
        <GiroRecep>Agricola</GiroRecep>
        <DirRecep>AVDA. VITACURA 2771 OF 1201</DirRecep>
        <CmnaRecep>Las Condes</CmnaRecep>
        <CiudadRecep>Santiago</CiudadRecep>
      </Receptor>
      <Totales>
        <MntNeto>35197</MntNeto>
        <TasaIVA>19</TasaIVA>
        <IVA>6687</IVA>
        <MntTotal>41884</MntTotal>
      </Totales>
    </Encabezado>
    <Detalle>
      <NroLinDet>1</NroLinDet>
      <CdgItem>
        <TpoCodigo>Interna</TpoCodigo>
        <VlrCodigo>EMP21</VlrCodigo>
      </CdgItem>
      <NmbItem>Plan Emprendedor </NmbItem>
      <DscItem>Abril 2025</DscItem>
      <QtyItem>0.90</QtyItem>
      <UnmdItem>Unid</UnmdItem>
      <PrcItem>39107.900000</PrcItem>
      <MontoItem>35197</MontoItem>
    </Detalle>
  </Documento>
</DTE>`)

	dte, err := ParseDTEXML(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse DTE XML: %v", err)
	}

	if dte.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", dte.Version)
	}

	doc := dte.Document
	if doc.ID != "DOC_29_33_2404" {
		t.Errorf("Expected document ID 'DOC_29_33_2404', got '%s'", doc.ID)
	}

	if doc.Header.DocInfo.DocumentType != 33 {
		t.Errorf("Expected document type 33, got %d", doc.Header.DocInfo.DocumentType)
	}

	if doc.Header.DocInfo.Folio != 2404 {
		t.Errorf("Expected folio 2404, got %d", doc.Header.DocInfo.Folio)
	}

	if doc.Header.DocInfo.IssueDate != "2025-05-05" {
		t.Errorf("Expected issue date '2025-05-05', got '%s'", doc.Header.DocInfo.IssueDate)
	}

	if doc.Header.Issuer.RUT != "76212889-6" {
		t.Errorf("Expected issuer RUT '76212889-6', got '%s'", doc.Header.Issuer.RUT)
	}

	if doc.Header.Issuer.CompanyName != "FACTURA MOVIL SPA" {
		t.Errorf("Expected issuer name 'FACTURA MOVIL SPA', got '%s'", doc.Header.Issuer.CompanyName)
	}

	if doc.Header.Receiver.RUT != "77371419-3" {
		t.Errorf("Expected receiver RUT '77371419-3', got '%s'", doc.Header.Receiver.RUT)
	}

	if doc.Header.Receiver.CompanyName != "AGRICOLA PAINE LTDA" {
		t.Errorf("Expected receiver name 'AGRICOLA PAINE LTDA', got '%s'", doc.Header.Receiver.CompanyName)
	}

	if doc.Header.Totals.NetAmount != 35197 {
		t.Errorf("Expected net amount 35197, got %f", doc.Header.Totals.NetAmount)
	}

	if doc.Header.Totals.TaxAmount != 6687 {
		t.Errorf("Expected tax amount 6687, got %f", doc.Header.Totals.TaxAmount)
	}

	if doc.Header.Totals.Total != 41884 {
		t.Errorf("Expected total 41884, got %f", doc.Header.Totals.Total)
	}

	if len(doc.Details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(doc.Details))
	}

	detail := doc.Details[0]
	if detail.LineNumber != 1 {
		t.Errorf("Expected line number 1, got %d", detail.LineNumber)
	}

	if detail.ItemName != "Plan Emprendedor " {
		t.Errorf("Expected item name 'Plan Emprendedor ', got '%s'", detail.ItemName)
	}

	if detail.Description != "Abril 2025" {
		t.Errorf("Expected description 'Abril 2025', got '%s'", detail.Description)
	}

	if detail.Quantity != 0.90 {
		t.Errorf("Expected quantity 0.90, got %f", detail.Quantity)
	}

	if detail.UnitPrice != 39107.900000 {
		t.Errorf("Expected unit price 39107.900000, got %f", detail.UnitPrice)
	}

	if detail.LineTotal != 35197 {
		t.Errorf("Expected line total 35197, got %f", detail.LineTotal)
	}
}

func TestDTEToInvoice(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="ISO-8859-1" standalone="no"?>
<DTE version="1.0">
  <Documento ID="DOC_29_33_2404">
    <Encabezado>
      <IdDoc>
        <TipoDTE>33</TipoDTE>
        <Folio>2404</Folio>
        <FchEmis>2025-05-05</FchEmis>
        <FmaPago>2</FmaPago>
        <FchVenc>2025-05-31</FchVenc>
      </IdDoc>
      <Emisor>
        <RUTEmisor>76212889-6</RUTEmisor>
        <RznSoc>FACTURA MOVIL SPA</RznSoc>
        <GiroEmis>COMERCIO AL POR MENOR DE COMPUTADORAS, SOFTWARES Y SUMINISTROS</GiroEmis>
        <CorreoEmisor>rfernandez@facturamovil.cl</CorreoEmisor>
        <Acteco>523930</Acteco>
        <Acteco>726000</Acteco>
        <DirOrigen>Vicuña Mackenna 9705</DirOrigen>
        <CmnaOrigen>La Florida</CmnaOrigen>
        <CiudadOrigen>Santiago</CiudadOrigen>
      </Emisor>
      <Receptor>
        <RUTRecep>77371419-3</RUTRecep>
        <RznSocRecep>AGRICOLA PAINE LTDA</RznSocRecep>
        <GiroRecep>Agricola</GiroRecep>
        <DirRecep>AVDA. VITACURA 2771 OF 1201</DirRecep>
        <CmnaRecep>Las Condes</CmnaRecep>
        <CiudadRecep>Santiago</CiudadRecep>
      </Receptor>
      <Totales>
        <MntNeto>35197</MntNeto>
        <TasaIVA>19</TasaIVA>
        <IVA>6687</IVA>
        <MntTotal>41884</MntTotal>
      </Totales>
    </Encabezado>
    <Detalle>
      <NroLinDet>1</NroLinDet>
      <CdgItem>
        <TpoCodigo>Interna</TpoCodigo>
        <VlrCodigo>EMP21</VlrCodigo>
      </CdgItem>
      <NmbItem>Plan Emprendedor </NmbItem>
      <DscItem>Abril 2025</DscItem>
      <QtyItem>0.90</QtyItem>
      <UnmdItem>Unid</UnmdItem>
      <PrcItem>39107.900000</PrcItem>
      <MontoItem>35197</MontoItem>
    </Detalle>
  </Documento>
</DTE>`)

	dte, err := ParseDTEXML(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse DTE XML: %v", err)
	}

	invoice, err := dte.ToInvoice()
	if err != nil {
		t.Fatalf("Failed to convert DTE to invoice: %v", err)
	}

	if invoice.DocumentType != 33 {
		t.Errorf("Expected document type 33, got %d", invoice.DocumentType)
	}

	if invoice.Folio != 2404 {
		t.Errorf("Expected folio 2404, got %d", invoice.Folio)
	}

	expectedDate := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	if !invoice.IssueDate.Equal(expectedDate) {
		t.Errorf("Expected issue date %v, got %v", expectedDate, invoice.IssueDate)
	}

	if invoice.Issuer.Code != "76212889-6" {
		t.Errorf("Expected issuer code '76212889-6', got '%s'", invoice.Issuer.Code)
	}

	if invoice.Issuer.Name != "FACTURA MOVIL SPA" {
		t.Errorf("Expected issuer name 'FACTURA MOVIL SPA', got '%s'", invoice.Issuer.Name)
	}

	if !strings.Contains(invoice.Issuer.Address, "Mackenna 9705") {
		t.Errorf("Expected issuer address to contain 'Mackenna 9705', got '%s'", invoice.Issuer.Address)
	}

	if invoice.Receiver == nil {
		t.Fatal("Expected receiver to be set")
	}

	if invoice.Receiver.Code != "77371419-3" {
		t.Errorf("Expected receiver code '77371419-3', got '%s'", invoice.Receiver.Code)
	}

	if invoice.Receiver.Name != "AGRICOLA PAINE LTDA" {
		t.Errorf("Expected receiver name 'AGRICOLA PAINE LTDA', got '%s'", invoice.Receiver.Name)
	}

	expectedReceiverAddress := "AVDA. VITACURA 2771 OF 1201, Las Condes, Santiago"
	if invoice.Receiver.Address != expectedReceiverAddress {
		t.Errorf("Expected receiver address '%s', got '%s'", expectedReceiverAddress, invoice.Receiver.Address)
	}

	if len(invoice.Details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(invoice.Details))
	}

	detail := invoice.Details[0]
	expectedDescription := "Plan Emprendedor  - Abril 2025"
	if detail.Description != expectedDescription {
		t.Errorf("Expected description '%s', got '%s'", expectedDescription, detail.Description)
	}

	if detail.Quantity != 0.90 {
		t.Errorf("Expected quantity 0.90, got %f", detail.Quantity)
	}

	if detail.UnitPrice != 39107.900000 {
		t.Errorf("Expected unit price 39107.900000, got %f", detail.UnitPrice)
	}

	if detail.LineTotal != 35197 {
		t.Errorf("Expected line total 35197, got %f", detail.LineTotal)
	}

	if invoice.Totals.TaxableAmount != 35197 {
		t.Errorf("Expected taxable amount 35197, got %f", invoice.Totals.TaxableAmount)
	}

	if invoice.Totals.TaxAmount != 6687 {
		t.Errorf("Expected tax amount 6687, got %f", invoice.Totals.TaxAmount)
	}

	if invoice.Totals.TotalAmount != 41884 {
		t.Errorf("Expected total amount 41884, got %f", invoice.Totals.TotalAmount)
	}

	total := invoice.CalculateTotal()
	if total != 41884 {
		t.Errorf("Expected calculated total 41884, got %d", total)
	}
}
