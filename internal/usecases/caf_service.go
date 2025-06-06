package usecases

import (
	"bytes"
	"context"
	"factura-movil-gateway/internal/domain"
	"fmt"
	"io"
)

// BlobStorageClient define la interfaz para almacenamiento de blobs.
type BlobStorageClient interface {
	Upload(ctx context.Context, blobName string, data io.Reader) error
}

// CAFRepository define la interfaz para gestionar CAFs (por ejemplo, guardar metadatos en BD).
type CAFRepository interface {
	Save(ctx context.Context, caf domain.CAF) error
	Update(ctx context.Context, caf domain.CAF) error
	FindByCompanyID(ctx context.Context, companyID string) ([]domain.CAF, error)
	FindAvailableCAF(ctx context.Context, companyID string, documentType uint) (*domain.CAF, error)
}

type CAFService interface {
	Create(ctx context.Context, company domain.Company, caf domain.CAF) error
	FindByCompanyID(ctx context.Context, companyID string) ([]domain.CAF, error)
	UseCAFFolio(ctx context.Context, companyID string, documentType uint) (int64, error)
}

func NewCAFService(storage BlobStorageClient, repository CAFRepository) *SimpleCAFService {
	return &SimpleCAFService{
		storage:    storage,
		repository: repository,
	}
}

type SimpleCAFService struct {
	storage    BlobStorageClient
	repository CAFRepository
}

func (s *SimpleCAFService) Create(ctx context.Context, company domain.Company, caf domain.CAF) error {
	err := s.repository.Save(ctx, caf)
	if err != nil {
		return fmt.Errorf("saving caf to database: %w", err)
	}

	blobName := fmt.Sprintf("caf/%s/%s.xml", company.ID, caf.ID)
	err = s.storage.Upload(ctx, blobName, bytes.NewReader(caf.Raw))
	if err != nil {
		return fmt.Errorf("uploading caf to storage: %w", err)
	}
	return nil
}

func (s *SimpleCAFService) FindByCompanyID(ctx context.Context, companyID string) ([]domain.CAF, error) {
	cafs, err := s.repository.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("finding cafs by company id: %w", err)
	}

	return cafs, nil
}

func (s *SimpleCAFService) UseCAFFolio(ctx context.Context, companyID string, documentType uint) (int64, error) {
	// Find an available CAF for this company and document type
	caf, err := s.repository.FindAvailableCAF(ctx, companyID, documentType)
	if err != nil {
		return 0, fmt.Errorf("finding available CAF: %w", err)
	}

	// Use the next folio
	folioToUse, shouldClose := caf.UseNextFolio()

	// Update the CAF in the database
	err = s.repository.Update(ctx, *caf)
	if err != nil {
		return 0, fmt.Errorf("updating CAF after folio use: %w", err)
	}

	if shouldClose {
		// Log that the CAF has been closed
		fmt.Printf("CAF %s has been closed after using all folios (final folio: %d)\n", caf.ID, caf.FinalFolios)
	}

	return folioToUse, nil
}
