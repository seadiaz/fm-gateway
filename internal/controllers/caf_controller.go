package controllers

import (
	"bytes"
	"encoding/xml"
	"factura-movil-gateway/internal/datatypes"
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/usecases"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	createCAFError           = "failed to create CAF"
	_cafCompanyNotFoundError = "company not found"

	_sixMonths = time.Hour * 24 * 30 * 6
)

func NewCAFController(cafService usecases.CAFService, companyService usecases.CompanyService) *CAFController {
	return &CAFController{
		cafService:     cafService,
		companyService: companyService,
	}
}

type CAFController struct {
	cafService     usecases.CAFService
	companyService usecases.CompanyService
}

func (c *CAFController) AddRoutes(mux *http.ServeMux) {
	mux.Handle("POST /companies/{companyId}/cafs", c.create())
	mux.Handle("GET /companies/{companyId}/cafs", c.list())
}

// decodeISO88591XML decodifica XML ISO-8859-1 a UTF-8.
func decodeISO88591XML(data []byte, v interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}

func (c *CAFController) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		companyId := r.PathValue("companyId")
		if companyId == "" {
			slog.Error("company id is required")
			httpserver.ReplyWithError(w, http.StatusBadRequest, _cafCompanyNotFoundError)
			return
		}

		company, err := c.companyService.FindByID(r.Context(), companyId)
		if err != nil {
			slog.Error("failed to find company", slog.String("Error", err.Error()), slog.String("companyId", companyId))
			httpserver.ReplyWithError(w, http.StatusNotFound, _cafCompanyNotFoundError)
			return
		}

		var body cafXML
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read request body", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, createCAFError)
			return
		}

		// Usa la función de decodificación personalizada
		err = decodeISO88591XML(rawData, &body)
		if err != nil {
			slog.Error("failed to decode xml", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, createCAFError)
			return
		}

		preview := body.RSASK.Value
		if len(preview) > 100 {
			preview = preview[:100]
		}
		slog.Debug("extracted private key from CAF",
			slog.Int("keyLength", len(body.RSASK.Value)),
			slog.String("keyPreview", preview))

		caf, err := domain.NewCAFBuilder().
			WithRaw(rawData).
			WithCompanyID(company.ID).
			WithCompanyCode(body.CAF.DA.RE).
			WithCompanyName(body.CAF.DA.RS).
			WithDocumentType(body.CAF.DA.TD).
			WithInitialFolios(body.CAF.DA.RNG.D).
			WithFinalFolios(body.CAF.DA.RNG.H).
			WithAuthorizationDate(body.CAF.DA.FA.Time).
			WithSignature(body.CAF.FRMA.Value).
			WithRSAPK_M(body.CAF.DA.RSAPK.M).
			WithRSAPK_E(body.CAF.DA.RSAPK.E).
			WithIDK(body.CAF.DA.IDK).
			WithPrivateKey(body.RSASK.Value).
			Build()
		if err != nil {
			slog.Error("failed to build CAF", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, createCAFError)
			return
		}

		err = c.cafService.Create(r.Context(), *company, caf)
		if err != nil {
			slog.Error("failed to create CAF", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, createCAFError)
			return
		}

		httpserver.ReplyJSONResponse(w, http.StatusCreated, caf)
	}
}

func (c *CAFController) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		companyId := r.PathValue("companyId")
		if companyId == "" {
			slog.Error("company id is required")
			httpserver.ReplyWithError(w, http.StatusBadRequest, _cafCompanyNotFoundError)
			return
		}

		// Verificar que la compañía existe
		_, err := c.companyService.FindByID(r.Context(), companyId)
		if err != nil {
			slog.Error("failed to find company", slog.String("Error", err.Error()), slog.String("companyId", companyId))
			httpserver.ReplyWithError(w, http.StatusNotFound, _cafCompanyNotFoundError)
			return
		}

		cafs, err := c.cafService.FindByCompanyID(r.Context(), companyId)
		if err != nil {
			slog.Error("failed to find cafs by company", slog.String("Error", err.Error()), slog.String("companyId", companyId))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, "failed to list cafs")
			return
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, cafs)
	}
}

// cafXML es una estructura interna para mapear el XML.
type cafXML struct {
	XMLName xml.Name `xml:"AUTORIZACION"`
	CAF     struct {
		DA struct {
			RE  string `xml:"RE"`
			RS  string `xml:"RS"`
			TD  uint   `xml:"TD"`
			RNG struct {
				D int64 `xml:"D"`
				H int64 `xml:"H"`
			} `xml:"RNG"`
			FA    datatypes.Date `xml:"FA"`
			RSAPK struct {
				M string `xml:"M"`
				E string `xml:"E"`
			} `xml:"RSAPK"`
			IDK string `xml:"IDK"`
		} `xml:"DA"`
		FRMA struct {
			Algorithm string `xml:"algoritmo,attr"`
			Value     string `xml:",chardata"`
		} `xml:"FRMA"`
	} `xml:"CAF"`
	RSASK struct {
		Value string `xml:",chardata"`
	} `xml:"RSASK"`
}
