package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"fmt"
	"time"
)

type StampService interface {
	Generate(ctx context.Context, company domain.Company, invoice domain.Invoice) (domain.Stamp, error)
}

func NewStampService(cafService CAFService) *SimpleStampService {
	return &SimpleStampService{
		cafService: cafService,
	}
}

type SimpleStampService struct {
	cafService CAFService
}

func (s *SimpleStampService) Generate(ctx context.Context, company domain.Company, invoice domain.Invoice) (domain.Stamp, error) {
	// Get the next available folio from CAF
	folio, err := s.cafService.UseCAFFolio(ctx, company.ID, uint(invoice.DocumentType))
	if err != nil {
		return domain.Stamp{}, fmt.Errorf("getting next folio from CAF: %w", err)
	}

	result := domain.Stamp{
		DD: domain.DD{
			RE:    company.Code,
			TD:    invoice.DocumentType,
			F:     folio,
			FE:    invoice.CreationDate,
			RR:    invoice.Customer.Code,
			RSR:   invoice.Customer.Name,
			MNT:   invoice.CalculateTotal(),
			IT1:   invoice.Details[0].Product.Name,
			CAF:   "FIXME",
			TSTED: time.Now().Format("2006-01-02T15:04:05-07:00"),
		},
	}

	return result, nil
}
