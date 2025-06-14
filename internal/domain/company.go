package domain

import "github.com/google/uuid"

// CommercialActivity represents a giro comercial (commercial activity)
type CommercialActivity struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"not null"`
}

// CompanyCommercialActivity represents the many-to-many relationship between companies and commercial activities
type CompanyCommercialActivity struct {
	CompanyID            string             `json:"company_id" gorm:"primaryKey"`
	CommercialActivityID string             `json:"commercial_activity_id" gorm:"primaryKey"`
	Company              Company            `json:"company" gorm:"foreignKey:CompanyID"`
	CommercialActivity   CommercialActivity `json:"commercial_activity" gorm:"foreignKey:CommercialActivityID"`
}

type Company struct {
	ID                    string               `json:"id" gorm:"primaryKey"`
	Code                  string               `json:"code" gorm:"uniqueIndex;not null"`
	Name                  string               `json:"name" gorm:"not null"`
	Address               string               `json:"address"`
	FacturaMovilCompanyID uint64               `json:"factura_movil_company_id"`
	CommercialActivities  []CommercialActivity `json:"commercial_activities" gorm:"many2many:company_commercial_activities"`
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

func (b *companyBuilder) WithCommercialActivities(activities []CommercialActivity) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.CommercialActivities = activities
		return nil
	})
	return b
}

func (b *companyBuilder) AddCommercialActivity(activity CommercialActivity) *companyBuilder {
	b.actions = append(b.actions, func(d *Company) error {
		d.CommercialActivities = append(d.CommercialActivities, activity)
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
