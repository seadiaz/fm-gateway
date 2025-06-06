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

type CompanyRepository interface {
	Save(ctx context.Context, company domain.Company) error
	FindAll(ctx context.Context) ([]domain.Company, error)
	FindByNameFilter(ctx context.Context, nameFilter string) ([]domain.Company, error)
	FindByID(ctx context.Context, id string) (*domain.Company, error)
}
