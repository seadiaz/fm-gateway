package controllers

import (
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/usecases"
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
				RE:    stamp.DD.RE,
				TD:    stamp.DD.TD,
				F:     stamp.DD.F,
				FE:    stamp.DD.FE,
				RR:    stamp.DD.RR,
				RSR:   stamp.DD.RSR,
				MNT:   stamp.DD.MNT,
				IT1:   stamp.DD.IT1,
				CAF:   stamp.DD.CAF,
				TSTED: stamp.DD.TSTED,
			},
			FRMT: stamp.FRMT,
		}

		httpserver.ReplyXMLResponse(w, http.StatusOK, response)
	}
}

// Estructuras para el JSON recibido
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
	FRMT    string `xml:"FRMT"`
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
	CAF   string `xml:"CAF"`
	TSTED string `xml:"TSTED"`
}
