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

	if err := db.AutoMigrate(&CompanyData{}, &CompanyCommercialActivityData{}); err != nil {
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
		Address:               company.Address,
		FacturaMovilCompanyID: company.FacturaMovilCompanyID,
	}
	err := c.db.
		WithContext(ctx).
		Save(&data).
		Error

	if err != nil {
		return fmt.Errorf("saving company: %w", err)
	}

	// Save commercial activities
	if len(company.CommercialActivities) > 0 {
		// First, remove existing activities
		err = c.db.WithContext(ctx).
			Where("company_id = ?", company.ID).
			Delete(&CompanyCommercialActivityData{}).Error
		if err != nil {
			return fmt.Errorf("removing existing commercial activities: %w", err)
		}

		// Then add new activities
		for _, activity := range company.CommercialActivities {
			activityData := CompanyCommercialActivityData{
				CompanyID:   company.ID,
				Code:        activity.Code,
				Description: activity.Description,
			}
			err = c.db.WithContext(ctx).Create(&activityData).Error
			if err != nil {
				return fmt.Errorf("saving commercial activity: %w", err)
			}
		}
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
		activities, err := c.GetCommercialActivities(ctx, data.ID)
		if err != nil {
			return nil, fmt.Errorf("getting commercial activities for company %s: %w", data.ID, err)
		}

		companies[i] = domain.Company{
			ID:                    data.ID,
			Name:                  data.Name,
			Code:                  data.Code,
			Address:               data.Address,
			FacturaMovilCompanyID: data.FacturaMovilCompanyID,
			CommercialActivities:  activities,
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
		activities, err := c.GetCommercialActivities(ctx, data.ID)
		if err != nil {
			return nil, fmt.Errorf("getting commercial activities for company %s: %w", data.ID, err)
		}

		companies[i] = domain.Company{
			ID:                    data.ID,
			Name:                  data.Name,
			Code:                  data.Code,
			Address:               data.Address,
			FacturaMovilCompanyID: data.FacturaMovilCompanyID,
			CommercialActivities:  activities,
		}
	}

	return companies, nil
}

func (c *CompanyRepository) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	if c.db == nil {
		return nil, errors.New("database not initialized")
	}

	var companyData CompanyData
	err := c.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&companyData).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("company not found with id: %s", id)
		}
		return nil, fmt.Errorf("finding company by id: %w", err)
	}

	activities, err := c.GetCommercialActivities(ctx, companyData.ID)
	if err != nil {
		return nil, fmt.Errorf("getting commercial activities for company %s: %w", companyData.ID, err)
	}

	company := domain.Company{
		ID:                    companyData.ID,
		Name:                  companyData.Name,
		Code:                  companyData.Code,
		Address:               companyData.Address,
		FacturaMovilCompanyID: companyData.FacturaMovilCompanyID,
		CommercialActivities:  activities,
	}

	return &company, nil
}

func (c *CompanyRepository) FindByCode(ctx context.Context, code string) (*domain.Company, error) {
	if c.db == nil {
		return nil, errors.New("database not initialized")
	}

	var companyData CompanyData
	err := c.db.
		WithContext(ctx).
		Where("code = ?", code).
		First(&companyData).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("company not found with code: %s", code)
		}
		return nil, fmt.Errorf("finding company by code: %w", err)
	}

	activities, err := c.GetCommercialActivities(ctx, companyData.ID)
	if err != nil {
		return nil, fmt.Errorf("getting commercial activities for company %s: %w", companyData.ID, err)
	}

	company := domain.Company{
		ID:                    companyData.ID,
		Name:                  companyData.Name,
		Code:                  companyData.Code,
		Address:               companyData.Address,
		FacturaMovilCompanyID: companyData.FacturaMovilCompanyID,
		CommercialActivities:  activities,
	}

	return &company, nil
}

func (c *CompanyRepository) GetCommercialActivities(ctx context.Context, companyID string) ([]domain.CommercialActivity, error) {
	if c.db == nil {
		return nil, errors.New("database not initialized")
	}

	var activitiesData []CompanyCommercialActivityData
	err := c.db.
		WithContext(ctx).
		Where("company_id = ?", companyID).
		Find(&activitiesData).
		Error

	if err != nil {
		return nil, fmt.Errorf("finding commercial activities: %w", err)
	}

	activities := make([]domain.CommercialActivity, len(activitiesData))
	for i, data := range activitiesData {
		activities[i] = domain.CommercialActivity{
			ID:          data.ID,
			Code:        data.Code,
			Description: data.Description,
		}
	}

	return activities, nil
}

func (c *CompanyRepository) AddCommercialActivity(ctx context.Context, companyID string, activity domain.CommercialActivity) error {
	if c.db == nil {
		return errors.New("database not initialized")
	}

	// First verify the company exists
	var company CompanyData
	err := c.db.WithContext(ctx).Where("id = ?", companyID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("company not found with id: %s", companyID)
		}
		return fmt.Errorf("finding company: %w", err)
	}

	// Create the activity
	activityData := CompanyCommercialActivityData{
		ID:          activity.ID,
		CompanyID:   companyID,
		Code:        activity.Code,
		Description: activity.Description,
	}

	err = c.db.WithContext(ctx).Create(&activityData).Error
	if err != nil {
		return fmt.Errorf("creating commercial activity: %w", err)
	}

	return nil
}

func (c *CompanyRepository) RemoveCommercialActivity(ctx context.Context, companyID string, activityID string) error {
	if c.db == nil {
		return errors.New("database not initialized")
	}

	// First verify the company exists
	var company CompanyData
	err := c.db.WithContext(ctx).Where("id = ?", companyID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("company not found with id: %s", companyID)
		}
		return fmt.Errorf("finding company: %w", err)
	}

	// Delete the activity
	err = c.db.WithContext(ctx).
		Where("company_id = ? AND id = ?", companyID, activityID).
		Delete(&CompanyCommercialActivityData{}).Error
	if err != nil {
		return fmt.Errorf("deleting commercial activity: %w", err)
	}

	return nil
}

type CompanyData struct {
	ID                    string `gorm:"primaryKey"`
	Name                  string
	Code                  string
	Address               string
	FacturaMovilCompanyID uint64
}

type CompanyCommercialActivityData struct {
	ID          string `gorm:"primaryKey"`
	CompanyID   string `gorm:"index;not null"`
	Code        string `gorm:"not null"`
	Description string `gorm:"not null"`
}
