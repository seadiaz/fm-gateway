package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
)

type CAFService interface {
	Create(ctx context.Context, caf domain.CAF) error
}

func NewCAFService(repository CAFRepository) *SimpleCAFService {
	return &SimpleCAFService{
		repository: repository,
	}
}

type SimpleCAFService struct {
	repository CAFRepository
}

func (s *SimpleCAFService) Create(ctx context.Context, caf domain.CAF) error {
	return s.repository.Save(ctx, caf)
}

// CAFRepository define la interfaz para gestionar CAFs.
type CAFRepository interface {
	Save(ctx context.Context, caf domain.CAF) error
}
