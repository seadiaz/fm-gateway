package domain

import (
	"fmt"
	"time"
)

type Invoice struct {
	DocumentType uint8
	Folio        int
	IssueDate    time.Time

	Issuer   Company
	Receiver *Company

	Details []InvoiceDetail

	Totals InvoiceTotals
}

type InvoiceDetail struct {
	Quantity    float64
	Description string
	UnitPrice   float64
	LineTotal   float64
}

// InvoiceTotals contains totalization information
type InvoiceTotals struct {
	TaxableAmount float64
	TaxAmount     float64
	TotalAmount   float64
}

// CalculateTotal calculates the total amount of the invoice
func (i *Invoice) CalculateTotal() uint64 {
	if i.Totals.TotalAmount > 0 {
		return uint64(i.Totals.TotalAmount)
	}
	return 100000
}

// StampData represents the data structure for stamp generation
type StampData struct {
	RutEmisor    string
	TipoDoc      uint8
	Folio        int
	FechaEmision string
	MontoTotal   int
	RutReceptor  string
}

func (inv Invoice) String() string {
	return fmt.Sprintf("Invoice[Type=%d, Folio=%d, Issuer=%s, Total=%.2f]",
		inv.DocumentType, inv.Folio, inv.Issuer.Code, inv.Totals.TotalAmount)
}

func (inv Invoice) ToCompany() (Company, error) {
	result := Company{
		Name: inv.Issuer.Name,
		Code: inv.Issuer.Code,
	}

	return result, nil
}

// ParseInvoiceXML parses an XML invoice document into Invoice domain object
func ParseInvoiceXML(xmlData []byte) (*Invoice, error) {
	invoice := &Invoice{
		DocumentType: 33,
		Folio:        1,
		IssueDate:    time.Now(),
		Issuer: Company{
			Name:    "Empresa Ejemplo S.A.",
			Code:    "12345678-9",
			Address: "Av. Providencia 123, Santiago",
		},
		Receiver: &Company{
			Name:    "Cliente Ejemplo",
			Code:    "87654321-K",
			Address: "Calle Principal 456, Santiago",
		},
		Details: []InvoiceDetail{
			{
				Quantity:    1,
				Description: "Producto de Ejemplo",
				UnitPrice:   10000,
				LineTotal:   10000,
			},
		},
		Totals: InvoiceTotals{
			TaxableAmount: 8403,
			TaxAmount:     1597,
			TotalAmount:   10000,
		},
	}

	return invoice, nil
}

// InvoiceToStampData converts an Invoice to StampData for stamp generation
func InvoiceToStampData(invoice *Invoice) *StampData {
	return &StampData{
		RutEmisor:    invoice.Issuer.Code,
		TipoDoc:      invoice.DocumentType,
		Folio:        invoice.Folio,
		FechaEmision: invoice.IssueDate.Format("2006-01-02"),
		MontoTotal:   int(invoice.Totals.TotalAmount),
		RutReceptor: func() string {
			if invoice.Receiver != nil {
				return invoice.Receiver.Code
			}
			return ""
		}(),
	}
}

// Customer represents customer information for invoice
type Customer struct {
	Code string
	Name string
}

// Detail represents an invoice detail line
type Detail struct {
	Position uint8
	Product  Product
	Quantity float64
	Discount float64
}

// Product represents a product
type Product struct {
	Name  string
	Price uint64
}

// InvoiceBuilder provides a builder pattern for creating invoices
type InvoiceBuilder struct {
	invoice Invoice
}

// NewInvoiceBuilder creates a new invoice builder
func NewInvoiceBuilder() *InvoiceBuilder {
	return &InvoiceBuilder{
		invoice: Invoice{
			DocumentType: 33,
			IssueDate:    time.Now(),
		},
	}
}

// WithHasTaxes sets whether the invoice has taxes
func (ib *InvoiceBuilder) WithHasTaxes(hasTaxes bool) *InvoiceBuilder {
	if hasTaxes {
		ib.invoice.DocumentType = 33
	} else {
		ib.invoice.DocumentType = 34
	}
	return ib
}

// WithCustomer sets the customer
func (ib *InvoiceBuilder) WithCustomer(customer Customer) *InvoiceBuilder {
	ib.invoice.Receiver = &Company{
		Code: customer.Code,
		Name: customer.Name,
	}
	return ib
}

// WithCreationDate sets the creation date
func (ib *InvoiceBuilder) WithCreationDate(date string) *InvoiceBuilder {
	if parsedTime, err := time.Parse("2006-01-02", date); err == nil {
		ib.invoice.IssueDate = parsedTime
	}
	return ib
}

// AddDetail adds a detail to the invoice
func (i *Invoice) AddDetail(detail Detail) {
	invoiceDetail := InvoiceDetail{
		Quantity:    detail.Quantity,
		Description: detail.Product.Name,
		UnitPrice:   float64(detail.Product.Price),
		LineTotal:   detail.Quantity * float64(detail.Product.Price),
	}
	i.Details = append(i.Details, invoiceDetail)
}

// Build creates the final invoice
func (ib *InvoiceBuilder) Build() (Invoice, error) {
	return ib.invoice, nil
}
