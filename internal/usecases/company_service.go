package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"fmt"
)

type CompanyService interface {
	Save(ctx context.Context, company domain.Company) error
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

type CompanyRepository interface {
	Save(ctx context.Context, company domain.Company) error
}
