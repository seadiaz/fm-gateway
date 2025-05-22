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
	createCAFError = "failed to create CAF"

	_sixMonths = time.Hour * 24 * 30 * 6
)

func NewCAFController(service usecases.CAFService) *CAFController {
	return &CAFController{
		service: service,
	}
}

type CAFController struct {
	service usecases.CAFService
}

func (c *CAFController) AddRoutes(mux *http.ServeMux) {
	mux.Handle("POST /caf", c.create())
}

// decodeISO88591XML decodifica XML ISO-8859-1 a UTF-8.
func decodeISO88591XML(data []byte, v interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}

func (c *CAFController) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body cafXML
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read request body", slog.String("Error", err.Error()))
			http.Error(w, createCAFError, http.StatusBadRequest)
			return
		}

		// Usa la función de decodificación personalizada
		err = decodeISO88591XML(rawData, &body)
		if err != nil {
			slog.Error("failed to decode xml", slog.String("Error", err.Error()))
			http.Error(w, createCAFError, http.StatusBadRequest)
			return
		}

		caf := domain.CAF{
			Raw:               rawData,
			CompanyID:         body.CAF.DA.RE,
			CompanyName:       body.CAF.DA.RS,
			DocumentType:      body.CAF.DA.TD,
			InitialFolios:     body.CAF.DA.RNG.D,
			FinalFolios:       body.CAF.DA.RNG.H,
			AuthorizationDate: body.CAF.DA.FA.Time,
			ExpirationDate:    body.CAF.DA.FA.Add(_sixMonths),
		}

		c.service.Create(r.Context(), caf)

		httpserver.ReplyJSONResponse(w, http.StatusCreated, caf)
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
			FA datatypes.Date `xml:"FA"`
		} `xml:"DA"`
	} `xml:"CAF"`
}
