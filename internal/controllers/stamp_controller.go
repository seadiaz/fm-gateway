package controllers

import (
	"encoding/base64"
	"encoding/xml"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/usecases"
	"factura-movil-gateway/internal/utils"
	"log/slog"
	"net/http"
)

const (
	_createStampError     = "failed to create stamp"
	_companyNotFoundError = "company not found"
)

func NewStampController(stampService usecases.StampService, companyService usecases.CompanyService) *StampController {
	return &StampController{
		stampService:   stampService,
		companyService: companyService,
	}
}

type StampController struct {
	stampService   usecases.StampService
	companyService usecases.CompanyService
}

func (c *StampController) AddRoutes(mux *http.ServeMux) {
	mux.Handle("POST /companies/{companyId}/stamps", c.create())
}

func (c *StampController) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		companyId := r.PathValue("companyId")
		if companyId == "" {
			slog.Error("company id is required")
			httpserver.ReplyWithError(w, http.StatusBadRequest, _companyNotFoundError)
			return
		}

		company, err := c.companyService.FindByID(r.Context(), companyId)
		if err != nil {
			slog.Error("failed to find company", slog.String("Error", err.Error()), slog.String("companyId", companyId))
			httpserver.ReplyWithError(w, http.StatusNotFound, _companyNotFoundError)
			return
		}

		var req StampRequest
		err = httpserver.DecodeJSONBody(r, &req)
		if err != nil {
			slog.Error("failed to decode json", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, _createStampError)
			return
		}

		invoice, err := domain.NewInvoiceBuilder().
			WithHasTaxes(req.HasTaxes).
			WithCustomer(domain.Customer{
				Code: req.Client.Code,
				Name: req.Client.Name,
			}).
			WithCreationDate(req.Date).
			Build()
		if err != nil {
			slog.Error("failed building invoice", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, _createStampError)
			return
		}

		for _, d := range req.Details {
			invoice.AddDetail(domain.Detail{
				Position: d.Position,
				Product: domain.Product{
					Name:  d.Product.Name,
					Price: d.Product.Price,
				},
				Quantity: d.Quantity,
				Discount: d.Discount,
			})
		}

		stamp, err := c.stampService.Generate(r.Context(), *company, invoice)
		if err != nil {
			slog.Error("failed to generate stamp", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _createStampError)
			return
		}

		response := TED{
			Version: "1.0",
			DD: DD{
				RE:  stamp.DD.RE,
				TD:  stamp.DD.TD,
				F:   stamp.DD.F,
				FE:  stamp.DD.FE,
				RR:  stamp.DD.RR,
				RSR: stamp.DD.RSR,
				MNT: stamp.DD.MNT,
				IT1: stamp.DD.IT1,
				CAF: CAF{
					Version: stamp.DD.CAF.Version,
					DA: DA{
						RE:  stamp.DD.CAF.DA.RE,
						RS:  stamp.DD.CAF.DA.RS,
						TD:  stamp.DD.CAF.DA.TD,
						RNG: RNG{D: stamp.DD.CAF.DA.RNG.D, H: stamp.DD.CAF.DA.RNG.H},
						FA:  stamp.DD.CAF.DA.FA,
						RSAPK: RSAPK{
							M: stamp.DD.CAF.DA.RSAPK.M,
							E: stamp.DD.CAF.DA.RSAPK.E,
						},
						IDK: stamp.DD.CAF.DA.IDK,
					},
					FRMA: FRMA{
						Algorithm: stamp.DD.CAF.FRMA.Algorithm,
						Value:     stamp.DD.CAF.FRMA.Value,
					},
				},
				TSTED: stamp.DD.TSTED,
			},
			FRMT: FRMT{
				Algorithm: "SHA1withRSA",
				Value:     stamp.FRMT,
			},
		}

		// Check if PDF417 barcode is requested via query parameter
		if r.URL.Query().Get("format") == "pdf417" || r.URL.Query().Get("include_barcode") == "true" {
			// Convert to XML for PDF417 generation
			xmlData, err := xml.Marshal(response)
			if err != nil {
				slog.Error("failed to marshal TED to XML for PDF417", slog.String("Error", err.Error()))
				httpserver.ReplyWithError(w, http.StatusInternalServerError, _createStampError)
				return
			}

			// Generate PDF417 barcode with automatically calculated safe dimensions
			_, pngBytes, err := utils.GenerateStampPDF417FromXMLWithAutoDimensions(string(xmlData))
			if err != nil {
				slog.Error("failed to generate PDF417 barcode", slog.String("Error", err.Error()))
				httpserver.ReplyWithError(w, http.StatusInternalServerError, _createStampError)
				return
			}

			// If only PDF417 is requested, return the image
			if r.URL.Query().Get("format") == "pdf417" {
				w.Header().Set("Content-Type", "image/png")
				w.WriteHeader(http.StatusOK)
				w.Write(pngBytes)
				return
			}

			// If include_barcode=true, return JSON with both XML and base64-encoded barcode
			barcodeResponse := StampWithBarcodeResponse{
				TED:     response,
				Barcode: base64.StdEncoding.EncodeToString(pngBytes),
			}
			httpserver.ReplyJSONResponse(w, http.StatusOK, barcodeResponse)
			return
		}

		// Default: return XML response
		httpserver.ReplyXMLResponse(w, http.StatusOK, response)
	}
}

type StampRequest struct {
	FmaPago       string     `json:"fmaPago"`
	HasTaxes      bool       `json:"hasTaxes"`
	Details       []Detail   `json:"details"`
	Client        Client     `json:"client"`
	AssignedFolio string     `json:"assignedFolio"`
	Subsidiary    Subsidiary `json:"subsidiary"`
	Date          string     `json:"date"`
}

type Detail struct {
	Position    uint8   `json:"position"`
	Product     Product `json:"product"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Discount    float64 `json:"discount"`
}

type Product struct {
	Unit  Unit   `json:"unit"`
	Price uint64 `json:"price"`
	Name  string `json:"name"`
	Code  string `json:"code"`
}

type Unit struct {
	Code string `json:"code"`
}

type Client struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	Municipality string `json:"municipality"`
	Line         string `json:"line"`
	Code         string `json:"code"`
}

type Subsidiary struct {
	Code string `json:"code"`
}

type TED struct {
	Version string `xml:"version,attr"`
	DD      DD     `xml:"DD"`
	FRMT    FRMT   `xml:"FRMT"`
}

type DD struct {
	RE    string `xml:"RE"`
	TD    uint8  `xml:"TD"`
	F     int64  `xml:"F"`
	FE    string `xml:"FE"`
	RR    string `xml:"RR"`
	RSR   string `xml:"RSR"`
	MNT   uint64 `xml:"MNT"`
	IT1   string `xml:"IT1"`
	CAF   CAF    `xml:"CAF"`
	TSTED string `xml:"TSTED"`
}

type FRMT struct {
	Algorithm string `xml:"algorithm,attr"`
	Value     string `xml:",chardata"`
}

type CAF struct {
	Version string `xml:"version,attr"`
	DA      DA     `xml:"DA"`
	FRMA    FRMA   `xml:"FRMA"`
}

type DA struct {
	RE    string `xml:"RE"`
	RS    string `xml:"RS"`
	TD    uint8  `xml:"TD"`
	RNG   RNG    `xml:"RNG"`
	FA    string `xml:"FA"`
	RSAPK RSAPK  `xml:"RSAPK"`
	IDK   string `xml:"IDK"`
}

type FRMA struct {
	Algorithm string `xml:"algorithm,attr"`
	Value     string `xml:",chardata"`
}

type RNG struct {
	D int64 `xml:"D"`
	H int64 `xml:"H"`
}

type RSAPK struct {
	M string `xml:"M"`
	E string `xml:"E"`
}

// StampWithBarcodeResponse includes both the TED XML structure and the PDF417 barcode
type StampWithBarcodeResponse struct {
	TED     TED    `json:"ted"`
	Barcode string `json:"barcode"` // Base64-encoded PNG image
}
