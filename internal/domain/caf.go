package domain

import "time"

type CAF struct {
	Raw               []byte
	CompanyID         string
	CompanyName       string
	DocumentType      uint
	InitialFolios     int64
	FinalFolios       int64
	AuthorizationDate time.Time
	ExpirationDate    time.Time
}
