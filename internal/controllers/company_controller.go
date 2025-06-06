package controllers

import (
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/usecases"
	"log/slog"
	"net/http"
)

const (
	_createCompanyError = "failed to create company"
	_listCompaniesError = "failed to list companies"
)

func NewCompanyController(service usecases.CompanyService) *CompanyController {
	return &CompanyController{
		service: service,
	}
}

type CompanyController struct {
	service usecases.CompanyService
}

func (c *CompanyController) AddRoutes(mux *http.ServeMux) {
	mux.Handle("POST /companies", c.create())
	mux.Handle("GET /companies", c.list())
}

func (c *CompanyController) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CompanyCreateRequest
		err := httpserver.DecodeJSONBody(r, &body)
		if err != nil {
			slog.Error("failed to decode json", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, _createCompanyError)
			return
		}

		company, err := domain.NewCompanyBuilder().
			WithName(body.Name).
			WithCode(body.Code).
			WithFacturaMovilCompanyID(body.FacturaMovilCompanyID).
			Build()
		if err != nil {
			slog.Error("failed to build domain company", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _createCompanyError)
			return
		}

		err = c.service.Save(r.Context(), company)
		if err != nil {
			slog.Error("failed to save company", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _createCompanyError)
			return
		}

		response := CompanyResponse{
			ID:                    company.ID,
			Name:                  company.Name,
			Code:                  company.Code,
			FacturaMovilCompanyID: company.FacturaMovilCompanyID,
		}

		httpserver.ReplyJSONResponse(w, http.StatusCreated, response)
	}
}

func (c *CompanyController) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameFilter := r.URL.Query().Get("name")

		var companies []domain.Company
		var err error

		if nameFilter != "" {
			companies, err = c.service.FindByNameFilter(r.Context(), nameFilter)
		} else {
			companies, err = c.service.FindAll(r.Context())
		}

		if err != nil {
			slog.Error("failed to find companies", slog.String("Error", err.Error()), slog.String("nameFilter", nameFilter))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _listCompaniesError)
			return
		}

		response := make([]CompanyResponse, len(companies))
		for i, company := range companies {
			response[i] = CompanyResponse{
				ID:                    company.ID,
				Name:                  company.Name,
				Code:                  company.Code,
				FacturaMovilCompanyID: company.FacturaMovilCompanyID,
			}
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, response)
	}
}

type CompanyCreateRequest struct {
	Name                  string `json:"name"`
	Code                  string `json:"code"`
	FacturaMovilCompanyID uint64 `json:"factura_movil_company_id"`
}

type CompanyResponse struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Code                  string `json:"code"`
	FacturaMovilCompanyID uint64 `json:"factura_movil_company_id"`
}
