package persistence

import (
	"bytes"
	"context"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/usecases"
	"fmt"
	"io"
	"time"
)

func NewCAFRepository(storage BlobStorageClient) *CAFRepository {
	return &CAFRepository{storage: storage}
}

var _ usecases.CAFRepository = (*CAFRepository)(nil)

type CAFRepository struct {
	storage BlobStorageClient
}

func (r *CAFRepository) Save(ctx context.Context, caf domain.CAF) error {
	blobName := fmt.Sprintf("caf/%d.xml", time.Now().UnixNano())
	err := r.storage.Upload(ctx, blobName, bytes.NewReader(caf.Raw))
	if err != nil {
		return err
	}

	return nil
}

type BlobStorageClient interface {
	Upload(ctx context.Context, blobName string, data io.Reader) error
}
