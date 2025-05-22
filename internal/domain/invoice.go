package domain

type Invoice struct {
	HasTaxes     bool
	Customer     Customer
	CreationDate string
	Details      []Detail

	DocumentType uint8
}

type Customer struct {
	Code string
	Name string
}

type Detail struct {
	Position uint8
	Product  Product
	Quantity float64
	Discount float64
}

type Product struct {
	Name  string
	Price uint64
}

func (i *Invoice) AddDetail(detail Detail) {
	i.Details = append(i.Details, detail)
}

func (i *Invoice) CalculateTotal() uint64 {
	var total uint64
	for _, d := range i.Details {
		total += uint64(d.Quantity * float64(d.Product.Price))
	}
	return total
}

func NewInvoiceBuilder() *invoiceBuilder {
	return &invoiceBuilder{}
}

type invoiceBuilder struct {
	actions []invoiceHandler
}

type invoiceHandler func(v *Invoice) error

func (b *invoiceBuilder) WithHasTaxes(value bool) *invoiceBuilder {
	b.actions = append(b.actions, func(d *Invoice) error {
		d.HasTaxes = value
		return nil
	})
	return b
}

func (b *invoiceBuilder) WithCustomer(value Customer) *invoiceBuilder {
	b.actions = append(b.actions, func(d *Invoice) error {
		d.Customer = value
		return nil
	})
	return b
}

func (b *invoiceBuilder) WithCreationDate(value string) *invoiceBuilder {
	b.actions = append(b.actions, func(d *Invoice) error {
		d.CreationDate = value
		return nil
	})
	return b
}

func (b *invoiceBuilder) Build() (Invoice, error) {
	result := Invoice{
		Details: make([]Detail, 0),
	}
	for _, a := range b.actions {
		if err := a(&result); err != nil {
			return Invoice{}, err
		}
	}

	if result.HasTaxes {
		result.DocumentType = 33
	} else {
		result.DocumentType = 34
	}

	return result, nil
}
