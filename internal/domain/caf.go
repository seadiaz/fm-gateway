package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	_sixMonths = time.Hour * 24 * 30 * 6
)

type CAF struct {
	ID                string
	Raw               []byte
	CompanyID         string
	CompanyCode       string
	CompanyName       string
	DocumentType      uint
	InitialFolios     int64
	CurrentFolios     int64
	FinalFolios       int64
	AuthorizationDate time.Time
	ExpirationDate    time.Time
}

func NewCAFBuilder() *cafBuilder {
	return &cafBuilder{}
}

type cafBuilder struct {
	actions []cafHandler
}

type cafHandler func(v *CAF) error

func (b *cafBuilder) WithRaw(value []byte) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.Raw = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithCompanyID(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.CompanyID = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithCompanyCode(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.CompanyCode = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithCompanyName(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.CompanyName = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithDocumentType(value uint) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.DocumentType = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithInitialFolios(value int64) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.InitialFolios = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithFinalFolios(value int64) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.FinalFolios = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithAuthorizationDate(value time.Time) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.AuthorizationDate = value
		return nil
	})
	return b
}

func (b *cafBuilder) Build() (CAF, error) {
	result := CAF{
		ID: uuid.NewString(),
	}
	for _, a := range b.actions {
		if err := a(&result); err != nil {
			return CAF{}, err
		}
	}

	result.CurrentFolios = result.InitialFolios
	result.ExpirationDate = result.AuthorizationDate.Add(_sixMonths)

	return result, nil
}
