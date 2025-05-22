package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implementa BlobStorageClient para almacenamiento local.
type LocalStorage struct {
	BasePath string
}

// NewLocalStorage crea una nueva instancia de LocalStorage.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

// Upload guarda el archivo localmente en la ruta base + blobName.
func (l *LocalStorage) Upload(ctx context.Context, blobName string, data io.Reader) error {
	fullPath := filepath.Join(l.BasePath, blobName)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	return err
}
