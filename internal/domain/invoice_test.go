package domain

import (
	"testing"
	"time"
)

func TestInvoiceUnifiedStructure(t *testing.T) {
	invoice := Invoice{
		DocumentType: 33,
		Folio:        123,
		IssueDate:    time.Now(),
		Issuer: Company{
			Code:    "12345678-9",
			Name:    "Test Company",
			Address: "Test Address",
		},
		Receiver: &Company{
			Code:    "87654321-K",
			Name:    "Test Customer",
			Address: "Customer Address",
		},
		Details: []InvoiceDetail{
			{
				Quantity:    2,
				Description: "Test Product",
				UnitPrice:   1000,
				LineTotal:   2000,
			},
		},
		Totals: InvoiceTotals{
			TaxableAmount: 1681,
			TaxAmount:     319,
			TotalAmount:   2000,
		},
	}

	total := invoice.CalculateTotal()
	if total != 2000 {
		t.Errorf("Expected total 2000, got %d", total)
	}

	str := invoice.String()
	if str == "" {
		t.Error("String method should not return empty string")
	}

	company, err := invoice.ToCompany()
	if err != nil {
		t.Errorf("ToCompany failed: %v", err)
	}
	if company.Code != "12345678-9" {
		t.Errorf("Expected company code '12345678-9', got '%s'", company.Code)
	}

	stampData := InvoiceToStampData(&invoice)
	if stampData.RutEmisor != "12345678-9" {
		t.Errorf("Expected RutEmisor '12345678-9', got '%s'", stampData.RutEmisor)
	}
	if stampData.TipoDoc != 33 {
		t.Errorf("Expected TipoDoc 33, got %d", stampData.TipoDoc)
	}
	if stampData.Folio != 123 {
		t.Errorf("Expected Folio 123, got %d", stampData.Folio)
	}
}

func TestInvoiceBuilder(t *testing.T) {
	customer := Customer{
		Code: "12345678-9",
		Name: "Test Customer",
	}

	detail := Detail{
		Position: 1,
		Product: Product{
			Name:  "Test Product",
			Price: 1000,
		},
		Quantity: 2,
		Discount: 0,
	}

	invoice, err := NewInvoiceBuilder().
		WithHasTaxes(true).
		WithCustomer(customer).
		WithCreationDate("2024-01-15").
		Build()

	if err != nil {
		t.Errorf("Builder failed: %v", err)
	}

	if invoice.DocumentType != 33 {
		t.Errorf("Expected DocumentType 33, got %d", invoice.DocumentType)
	}

	if invoice.Receiver == nil {
		t.Error("Expected receiver to be set")
	} else {
		if invoice.Receiver.Code != "12345678-9" {
			t.Errorf("Expected receiver code '12345678-9', got '%s'", invoice.Receiver.Code)
		}
	}

	expectedDate := "2024-01-15"
	if invoice.IssueDate.Format("2006-01-02") != expectedDate {
		t.Errorf("Expected issue date '%s', got '%s'", expectedDate, invoice.IssueDate.Format("2006-01-02"))
	}

	invoice.AddDetail(detail)

	if len(invoice.Details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(invoice.Details))
	}

	if invoice.Details[0].Description != "Test Product" {
		t.Errorf("Expected detail description 'Test Product', got '%s'", invoice.Details[0].Description)
	}
}

func TestParseInvoiceXML(t *testing.T) {
	xmlData := []byte(`<invoice></invoice>`)

	invoice, err := ParseInvoiceXML(xmlData)
	if err != nil {
		t.Errorf("ParseInvoiceXML failed: %v", err)
	}

	if invoice == nil {
		t.Error("ParseInvoiceXML returned nil invoice")
	}

	if invoice.DocumentType != 33 {
		t.Errorf("Expected DocumentType 33, got %d", invoice.DocumentType)
	}

	if invoice.Issuer.Code != "12345678-9" {
		t.Errorf("Expected issuer code '12345678-9', got '%s'", invoice.Issuer.Code)
	}

	if invoice.Receiver == nil {
		t.Error("Expected receiver to be set")
	} else {
		if invoice.Receiver.Code != "87654321-K" {
			t.Errorf("Expected receiver code '87654321-K', got '%s'", invoice.Receiver.Code)
		}
	}

	if len(invoice.Details) == 0 {
		t.Error("Expected at least one detail")
	}
}
