package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"strings"
	"testing"
	"time"
)

// Mock CompanyService for testing
type mockCompanyService struct {
	companies map[string]*domain.Company
}

func (m *mockCompanyService) Save(ctx context.Context, company domain.Company) error {
	return nil
}

func (m *mockCompanyService) FindAll(ctx context.Context) ([]domain.Company, error) {
	return nil, nil
}

func (m *mockCompanyService) FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error) {
	return nil, nil
}

func (m *mockCompanyService) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	return nil, nil
}

func (m *mockCompanyService) FindByCode(ctx context.Context, code string) (*domain.Company, error) {
	if company, exists := m.companies[code]; exists {
		return company, nil
	}
	return nil, &CompanyNotFoundError{Code: code}
}

// Mock StampService for testing
type mockStampService struct{}

func (m *mockStampService) Generate(ctx context.Context, company domain.Company, invoice domain.Invoice) (domain.Stamp, error) {
	return domain.Stamp{
		DD: domain.DD{
			RE:  company.Code,
			TD:  invoice.DocumentType,
			F:   int64(invoice.Folio),
			FE:  invoice.IssueDate.Format("2006-01-02"),
			RR:  invoice.Receiver.Code,
			RSR: invoice.Receiver.Name,
			MNT: uint64(invoice.Totals.TotalAmount),
		},
		FRMT: "mock-signature",
	}, nil
}

// Custom error type for company not found
type CompanyNotFoundError struct {
	Code string
}

func (e *CompanyNotFoundError) Error() string {
	return "company not found with code: " + e.Code
}

func TestDocumentService_ProcessInvoice_CompanyLookup(t *testing.T) {
	// Setup mock services
	testCompany := &domain.Company{
		ID:   "test-id",
		Code: "12345678-9",
		Name: "Test Company",
	}

	companyService := &mockCompanyService{
		companies: map[string]*domain.Company{
			"12345678-9": testCompany,
		},
	}

	stampService := &mockStampService{}
	documentService := NewDocumentService(stampService, companyService)

	// Create test invoice
	invoice := &domain.Invoice{
		DocumentType: 33,
		Folio:        123,
		IssueDate:    time.Now(),
		Issuer: domain.Company{
			Code: "12345678-9",
			Name: "Test Company",
		},
		Receiver: &domain.Company{
			Code: "87654321-0",
			Name: "Test Receiver",
		},
		Details: []domain.InvoiceDetail{
			{
				Description: "Test Item",
				Quantity:    1,
				UnitPrice:   1000,
				LineTotal:   1000,
			},
		},
		Totals: domain.InvoiceTotals{
			TaxableAmount: 840,
			TaxAmount:     160,
			TotalAmount:   1000,
		},
	}

	// Test processing
	result, err := documentService.ProcessInvoice(invoice)

	// Verify results
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.StampXML) == 0 {
		t.Error("Expected stamp XML to be generated")
	}

	if len(result.PDF417Data) == 0 {
		t.Error("Expected PDF417 data to be generated")
	}

	if len(result.ThermalPDF) == 0 {
		t.Error("Expected thermal PDF to be generated")
	}

	if result.ProcessingTime == 0 {
		t.Error("Expected processing time to be recorded")
	}
}

func TestDocumentService_ProcessInvoice_CompanyNotFound(t *testing.T) {
	// Setup mock services with empty company map
	companyService := &mockCompanyService{
		companies: map[string]*domain.Company{},
	}

	stampService := &mockStampService{}
	documentService := NewDocumentService(stampService, companyService)

	// Create test invoice with non-existent company
	invoice := &domain.Invoice{
		DocumentType: 33,
		Folio:        123,
		IssueDate:    time.Now(),
		Issuer: domain.Company{
			Code: "99999999-9", // Non-existent company
			Name: "Non-existent Company",
		},
		Receiver: &domain.Company{
			Code: "87654321-0",
			Name: "Test Receiver",
		},
		Details: []domain.InvoiceDetail{
			{
				Description: "Test Item",
				Quantity:    1,
				UnitPrice:   1000,
				LineTotal:   1000,
			},
		},
		Totals: domain.InvoiceTotals{
			TaxableAmount: 840,
			TaxAmount:     160,
			TotalAmount:   1000,
		},
	}

	// Test processing
	result, err := documentService.ProcessInvoice(invoice)

	// Verify error handling
	if err == nil {
		t.Fatal("Expected error for non-existent company, got nil")
	}

	if result.Error == nil {
		t.Error("Expected result.Error to be set")
	}

	// Check that error message contains company code
	if !strings.Contains(err.Error(), "99999999-9") {
		t.Errorf("Expected error message to contain company code, got: %v", err)
	}
}
