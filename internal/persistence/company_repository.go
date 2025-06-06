package persistence

import (
	"context"
	"errors"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/usecases"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewCompanyRepository(dsn string) (*CompanyRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&CompanyData{}); err != nil {
		return nil, err
	}
	return &CompanyRepository{db: db}, nil
}

var _ usecases.CompanyRepository = (*CompanyRepository)(nil)

type CompanyRepository struct {
	db *gorm.DB
}

func (c *CompanyRepository) Save(ctx context.Context, company domain.Company) error {
	if c.db == nil {
		return errors.New("database not initialized")
	}

	data := CompanyData{
		ID:                    company.ID,
		Name:                  company.Name,
		Code:                  company.Code,
		FacturaMovilCompanyID: company.FacturaMovilCompanyID,
	}
	err := c.db.
		WithContext(ctx).
		Create(&data).
		Error

	if err != nil {
		return fmt.Errorf("saving company: %w", err)
	}

	return nil
}

func (c *CompanyRepository) FindAll(ctx context.Context) ([]domain.Company, error) {
	if c.db == nil {
		return nil, errors.New("database not initialized")
	}

	var companiesData []CompanyData
	err := c.db.
		WithContext(ctx).
		Find(&companiesData).
		Error

	if err != nil {
		return nil, fmt.Errorf("finding all companies: %w", err)
	}

	companies := make([]domain.Company, len(companiesData))
	for i, data := range companiesData {
		companies[i] = domain.Company{
			ID:                    data.ID,
			Name:                  data.Name,
			Code:                  data.Code,
			FacturaMovilCompanyID: data.FacturaMovilCompanyID,
		}
	}

	return companies, nil
}

func (c *CompanyRepository) FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error) {
	if c.db == nil {
		return nil, errors.New("database not initialized")
	}

	var companiesData []CompanyData
	err := c.db.
		WithContext(ctx).
		Where("name ILIKE ?", "%"+nameFilter+"%").
		Find(&companiesData).
		Error

	if err != nil {
		return nil, fmt.Errorf("finding companies by name filter: %w", err)
	}

	companies := make([]domain.Company, len(companiesData))
	for i, data := range companiesData {
		companies[i] = domain.Company{
			ID:                    data.ID,
			Name:                  data.Name,
			Code:                  data.Code,
			FacturaMovilCompanyID: data.FacturaMovilCompanyID,
		}
	}

	return companies, nil
}

type CompanyData struct {
	ID                    string `gorm:"primaryKey"`
	Name                  string
	Code                  string
	FacturaMovilCompanyID uint64
}
