package domain

import "github.com/google/uuid"

type Company struct {
	ID                    string
	Code                  string
	Name                  string
	Address               string
	FacturaMovilCompanyID uint64
}

func NewCompanyBuilder() *companyBuilder {
	return &companyBuilder{}
}

type companyBuilder struct {
	actions []compnayHandler
}

type compnayHandler func(v *Company) error

func (b *companyBuilder) WithCode(value string) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.Code = value
		return nil
	})
	return b
}

func (b *companyBuilder) WithName(value string) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.Name = value
		return nil
	})
	return b
}

func (b *companyBuilder) WithAddress(value string) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.Address = value
		return nil
	})
	return b
}

func (b *companyBuilder) WithFacturaMovilCompanyID(value uint64) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.FacturaMovilCompanyID = value
		return nil
	})
	return b
}

func (b *companyBuilder) Build() (Company, error) {
	result := Company{
		ID: uuid.NewString(),
	}
	for _, a := range b.actions {
		if err := a(&result); err != nil {
			return Company{}, err
		}
	}

	return result, nil
}
