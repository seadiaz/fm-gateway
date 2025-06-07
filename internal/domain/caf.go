package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	_sixMonths = time.Hour * 24 * 30 * 6
)

// CAF status constants
const (
	CAFStatusOpen   = "OPEN"
	CAFStatusClosed = "CLOSED"
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
	Status            string
	Signature         string
	RSAPK_M           string
	RSAPK_E           string
	IDK               string
	PrivateKey        string
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

func (b *cafBuilder) WithSignature(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.Signature = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithRSAPK_M(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.RSAPK_M = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithRSAPK_E(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.RSAPK_E = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithIDK(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.IDK = value
		return nil
	})
	return b
}

func (b *cafBuilder) WithPrivateKey(value string) *cafBuilder {
	b.actions = append(b.actions, func(d *CAF) error {
		d.PrivateKey = value
		return nil
	})
	return b
}

func (b *cafBuilder) Build() (CAF, error) {
	result := CAF{
		ID:     uuid.NewString(),
		Status: CAFStatusOpen, // CAFs start as open
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

// UseNextFolio increments the current folio and returns the folio to use
// Returns the folio number and whether the CAF should be closed after this use
func (c *CAF) UseNextFolio() (int64, bool) {
	folioToUse := c.CurrentFolios
	shouldClose := c.CurrentFolios == c.FinalFolios
	c.CurrentFolios++

	if shouldClose {
		c.Status = CAFStatusClosed
	}

	return folioToUse, shouldClose
}

// IsOpen returns true if the CAF is in open status
func (c *CAF) IsOpen() bool {
	return c.Status == CAFStatusOpen
}

// HasAvailableFolios returns true if there are folios available to use
func (c *CAF) HasAvailableFolios() bool {
	return c.IsOpen() && c.CurrentFolios <= c.FinalFolios
}
