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
}

type CAFService interface {
	Create(ctx context.Context, caf domain.CAF) error
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

func (s *SimpleCAFService) Create(ctx context.Context, caf domain.CAF) error {
	err := s.repository.Save(ctx, caf)
	if err != nil {
		return fmt.Errorf("saving caf to database: %w", err)
	}

	blobName := fmt.Sprintf("caf/%s.xml", caf.ID)
	err = s.storage.Upload(ctx, blobName, bytes.NewReader(caf.Raw))
	if err != nil {
		return fmt.Errorf("uploading caf to storage: %w", err)
	}
	return nil
}
