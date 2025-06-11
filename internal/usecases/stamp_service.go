package usecases

import (
	"context"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/utils"
	"fmt"
	"time"
)

type StampService interface {
	Generate(ctx context.Context, company domain.Company, invoice domain.Invoice) (domain.Stamp, error)
}

func NewStampService(cafService CAFService) *SimpleStampService {
	return &SimpleStampService{
		cafService: cafService,
	}
}

type SimpleStampService struct {
	cafService CAFService
}

func (s *SimpleStampService) Generate(ctx context.Context, company domain.Company, invoice domain.Invoice) (domain.Stamp, error) {
	folio, caf, err := s.cafService.UseCAFFolio(ctx, company.ID, uint(invoice.DocumentType))
	if err != nil {
		return domain.Stamp{}, fmt.Errorf("getting next folio from CAF: %w", err)
	}

	// Create StampCAF from domain CAF
	stampCAF := domain.StampCAF{
		Version: "1.0",
		DA: domain.StampDA{
			RE: caf.CompanyCode,
			RS: caf.CompanyName,
			TD: uint8(caf.DocumentType),
			RNG: domain.StampRNG{
				D: caf.InitialFolios,
				H: caf.FinalFolios,
			},
			FA: caf.AuthorizationDate.Format("2006-01-02"),
			RSAPK: domain.StampRSAPK{
				M: caf.RSAPK_M,
				E: caf.RSAPK_E,
			},
			IDK: caf.IDK,
		},
		FRMA: domain.StampFRMA{
			Algorithm: "SHA1withRSA",
			Value:     caf.Signature,
		},
	}

	// Create DD structure
	dd := domain.DD{
		RE:    company.Code,
		TD:    invoice.DocumentType,
		F:     folio,
		FE:    invoice.CreationDate,
		RR:    invoice.Customer.Code,
		RSR:   invoice.Customer.Name,
		MNT:   invoice.CalculateTotal(),
		IT1:   invoice.Details[0].Product.Name,
		CAF:   stampCAF,
		TSTED: time.Now().Format("2006-01-02T15:04:05"),
	}

	// Serialize DD to XML without newlines
	ddXML, err := utils.SerializeToXMLWithoutNewlines(dd)
	if err != nil {
		return domain.Stamp{}, fmt.Errorf("serializing DD to XML: %w", err)
	}

	// Sign the DD XML with the CAF private key
	frmt, err := utils.SignSHA1WithRSA(ddXML, caf.PrivateKey)
	if err != nil {
		return domain.Stamp{}, fmt.Errorf("signing DD with private key (key length: %d): %w", len(caf.PrivateKey), err)
	}

	result := domain.Stamp{
		DD:   dd,
		FRMT: frmt,
	}

	return result, nil
}
