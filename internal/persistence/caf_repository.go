package persistence

import (
	"context"
	"errors"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/usecases"
	"fmt"
	"io"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewCAFRepository(dsn string) (*CAFRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&CAFData{}); err != nil {
		return nil, err
	}
	return &CAFRepository{db: db}, nil
}

var _ usecases.CAFRepository = (*CAFRepository)(nil)

type CAFRepository struct {
	db *gorm.DB
}

func (r *CAFRepository) Save(ctx context.Context, caf domain.CAF) error {
	if r.db == nil {
		return errors.New("database not initialized")
	}

	data := CAFData{
		ID:                caf.ID,
		Raw:               caf.Raw,
		CompanyID:         caf.CompanyID,
		CompanyCode:       caf.CompanyCode,
		CompanyName:       caf.CompanyName,
		DocumentType:      caf.DocumentType,
		InitialFolios:     caf.InitialFolios,
		CurrentFolios:     caf.CurrentFolios,
		FinalFolios:       caf.FinalFolios,
		AuthorizationDate: caf.AuthorizationDate,
		ExpirationDate:    caf.ExpirationDate,
	}
	err := r.db.
		WithContext(ctx).
		Create(&data).
		Error

	if err != nil {
		return fmt.Errorf("saving caf: %w", err)
	}

	return nil
}

func (r *CAFRepository) FindByCompanyID(ctx context.Context, companyID string) ([]domain.CAF, error) {
	if r.db == nil {
		return nil, errors.New("database not initialized")
	}

	var cafsData []CAFData
	err := r.db.
		WithContext(ctx).
		Where("company_id = ?", companyID).
		Find(&cafsData).
		Error

	if err != nil {
		return nil, fmt.Errorf("finding cafs by company id: %w", err)
	}

	cafs := make([]domain.CAF, len(cafsData))
	for i, data := range cafsData {
		cafs[i] = domain.CAF{
			ID:                data.ID,
			Raw:               data.Raw,
			CompanyID:         data.CompanyID,
			CompanyCode:       data.CompanyCode,
			CompanyName:       data.CompanyName,
			DocumentType:      data.DocumentType,
			InitialFolios:     data.InitialFolios,
			CurrentFolios:     data.CurrentFolios,
			FinalFolios:       data.FinalFolios,
			AuthorizationDate: data.AuthorizationDate,
			ExpirationDate:    data.ExpirationDate,
		}
	}

	return cafs, nil
}

type BlobStorageClient interface {
	Upload(ctx context.Context, blobName string, data io.Reader) error
}

type CAFData struct {
	ID                string `gorm:"primaryKey"`
	Raw               []byte
	CompanyID         string `gorm:"index"`
	CompanyCode       string
	CompanyName       string
	DocumentType      uint
	InitialFolios     int64
	CurrentFolios     int64
	FinalFolios       int64
	AuthorizationDate time.Time
	ExpirationDate    time.Time
}
