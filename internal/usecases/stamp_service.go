package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"time"
)

type StampService interface {
	Generate(ctx context.Context, invoice domain.Invoice) (domain.Stamp, error)
}

func NewStampService() *SimpleStampService {
	return &SimpleStampService{}
}

type SimpleStampService struct {
}

func (s *SimpleStampService) Generate(ctx context.Context, invoice domain.Invoice) (domain.Stamp, error) {
	result := domain.Stamp{
		DD: domain.DD{
			RE:    "12345678-9",
			TD:    invoice.DocumentType,
			F:     1,
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
