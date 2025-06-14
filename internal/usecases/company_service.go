package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"fmt"
)

type CompanyService interface {
	Save(ctx context.Context, company domain.Company) error
	FindAll(ctx context.Context) ([]domain.Company, error)
	FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error)
	FindByID(ctx context.Context, id string) (*domain.Company, error)
	FindByCode(ctx context.Context, code string) (*domain.Company, error)
	Update(ctx context.Context, company domain.Company) error
	AddCommercialActivity(ctx context.Context, companyID string, activity domain.CommercialActivity) error
	RemoveCommercialActivity(ctx context.Context, companyID string, activityID string) error
	GetCommercialActivities(ctx context.Context, companyID string) ([]domain.CommercialActivity, error)
}

func NewCompanyService(repository CompanyRepository) *SimpleCompanyService {
	return &SimpleCompanyService{
		reponsitory: repository,
	}
}

type SimpleCompanyService struct {
	reponsitory CompanyRepository
}

func (s *SimpleCompanyService) Save(ctx context.Context, company domain.Company) error {
	err := s.reponsitory.Save(ctx, company)
	if err != nil {
		return fmt.Errorf("saving company: %w", err)
	}

	return nil
}

func (s *SimpleCompanyService) FindAll(ctx context.Context) ([]domain.Company, error) {
	companies, err := s.reponsitory.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all companies: %w", err)
	}

	return companies, nil
}

func (s *SimpleCompanyService) FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error) {
	companies, err := s.reponsitory.FindByNameFilter(ctx, nameFilter)
	if err != nil {
		return nil, fmt.Errorf("finding companies by name filter: %w", err)
	}

	return companies, nil
}

func (s *SimpleCompanyService) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	company, err := s.reponsitory.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding company by id: %w", err)
	}

	return company, nil
}

func (s *SimpleCompanyService) FindByCode(ctx context.Context, code string) (*domain.Company, error) {
	company, err := s.reponsitory.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("finding company by code: %w", err)
	}

	return company, nil
}

func (s *SimpleCompanyService) Update(ctx context.Context, company domain.Company) error {
	err := s.reponsitory.Save(ctx, company)
	if err != nil {
		return fmt.Errorf("updating company: %w", err)
	}

	return nil
}

func (s *SimpleCompanyService) AddCommercialActivity(ctx context.Context, companyID string, activity domain.CommercialActivity) error {
	err := s.reponsitory.AddCommercialActivity(ctx, companyID, activity)
	if err != nil {
		return fmt.Errorf("adding commercial activity: %w", err)
	}

	return nil
}

func (s *SimpleCompanyService) RemoveCommercialActivity(ctx context.Context, companyID string, activityID string) error {
	err := s.reponsitory.RemoveCommercialActivity(ctx, companyID, activityID)
	if err != nil {
		return fmt.Errorf("removing commercial activity: %w", err)
	}

	return nil
}

func (s *SimpleCompanyService) GetCommercialActivities(ctx context.Context, companyID string) ([]domain.CommercialActivity, error) {
	activities, err := s.reponsitory.GetCommercialActivities(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("getting commercial activities: %w", err)
	}

	return activities, nil
}

type CompanyRepository interface {
	Save(ctx context.Context, company domain.Company) error
	FindAll(ctx context.Context) ([]domain.Company, error)
	FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error)
	FindByID(ctx context.Context, id string) (*domain.Company, error)
	FindByCode(ctx context.Context, code string) (*domain.Company, error)
	GetCommercialActivities(ctx context.Context, companyID string) ([]domain.CommercialActivity, error)
	AddCommercialActivity(ctx context.Context, companyID string, activity domain.CommercialActivity) error
	RemoveCommercialActivity(ctx context.Context, companyID string, activityID string) error
}
