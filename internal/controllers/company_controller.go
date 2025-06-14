package controllers

import (
	"factura-movil-gateway/internal/domain"
	"factura-movil-gateway/internal/httpserver"
	"factura-movil-gateway/internal/usecases"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	_createCompanyError            = "failed to create company"
	_listCompaniesError            = "failed to list companies"
	_getCompanyError               = "failed to get company"
	_updateCompanyError            = "failed to update company"
	_addCommercialActivityError    = "failed to add commercial activity"
	_removeCommercialActivityError = "failed to remove commercial activity"
	_getCommercialActivitiesError  = "failed to get commercial activities"
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
	mux.Handle("GET /companies/{id}", c.getByID())
	mux.Handle("PUT /companies/{id}", c.update())
	mux.Handle("POST /companies/{id}/commercial-activities", c.addCommercialActivity())
	mux.Handle("DELETE /companies/{id}/commercial-activities/{activityId}", c.removeCommercialActivity())
	mux.Handle("GET /companies/{id}/commercial-activities", c.getCommercialActivities())
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
			WithAddress(body.Address).
			WithFacturaMovilCompanyID(body.FacturaMovilCompanyID).
			WithCommercialActivities(body.CommercialActivities).
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
			Address:               company.Address,
			FacturaMovilCompanyID: company.FacturaMovilCompanyID,
			CommercialActivities:  make([]CommercialActivityResponse, len(company.CommercialActivities)),
		}

		for i, activity := range company.CommercialActivities {
			response.CommercialActivities[i] = CommercialActivityResponse{
				ID:          activity.ID,
				Code:        activity.Code,
				Description: activity.Description,
			}
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
				Address:               company.Address,
				FacturaMovilCompanyID: company.FacturaMovilCompanyID,
				CommercialActivities:  make([]CommercialActivityResponse, len(company.CommercialActivities)),
			}

			for j, activity := range company.CommercialActivities {
				response[i].CommercialActivities[j] = CommercialActivityResponse{
					ID:          activity.ID,
					Code:        activity.Code,
					Description: activity.Description,
				}
			}
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, response)
	}
}

func (c *CompanyController) getByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/companies/")
		if id == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "id is required")
			return
		}

		company, err := c.service.FindByID(r.Context(), id)
		if err != nil {
			slog.Error("failed to find company", slog.String("Error", err.Error()), slog.String("id", id))
			httpserver.ReplyWithError(w, http.StatusNotFound, _getCompanyError)
			return
		}

		response := CompanyResponse{
			ID:                    company.ID,
			Name:                  company.Name,
			Code:                  company.Code,
			Address:               company.Address,
			FacturaMovilCompanyID: company.FacturaMovilCompanyID,
			CommercialActivities:  make([]CommercialActivityResponse, len(company.CommercialActivities)),
		}

		for i, activity := range company.CommercialActivities {
			response.CommercialActivities[i] = CommercialActivityResponse{
				ID:          activity.ID,
				Code:        activity.Code,
				Description: activity.Description,
			}
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, response)
	}
}

func (c *CompanyController) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/companies/")
		if id == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "id is required")
			return
		}

		var body CompanyUpdateRequest
		err := httpserver.DecodeJSONBody(r, &body)
		if err != nil {
			slog.Error("failed to decode json", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, _updateCompanyError)
			return
		}

		company := domain.Company{
			ID:                    id,
			Name:                  body.Name,
			Code:                  body.Code,
			Address:               body.Address,
			FacturaMovilCompanyID: body.FacturaMovilCompanyID,
			CommercialActivities:  body.CommercialActivities,
		}

		err = c.service.Update(r.Context(), company)
		if err != nil {
			slog.Error("failed to update company", slog.String("Error", err.Error()), slog.String("id", id))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _updateCompanyError)
			return
		}

		response := CompanyResponse{
			ID:                    company.ID,
			Name:                  company.Name,
			Code:                  company.Code,
			Address:               company.Address,
			FacturaMovilCompanyID: company.FacturaMovilCompanyID,
			CommercialActivities:  make([]CommercialActivityResponse, len(company.CommercialActivities)),
		}

		for i, activity := range company.CommercialActivities {
			response.CommercialActivities[i] = CommercialActivityResponse{
				ID:          activity.ID,
				Code:        activity.Code,
				Description: activity.Description,
			}
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, response)
	}
}

func (c *CompanyController) addCommercialActivity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "company id is required")
			return
		}

		var body CommercialActivityCreateRequest
		err := httpserver.DecodeJSONBody(r, &body)
		if err != nil {
			slog.Error("failed to decode json", slog.String("Error", err.Error()))
			httpserver.ReplyWithError(w, http.StatusBadRequest, _addCommercialActivityError)
			return
		}

		if body.Code == "" || body.Description == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "code and description are required")
			return
		}

		activity := domain.CommercialActivity{
			ID:          uuid.NewString(),
			Code:        body.Code,
			Description: body.Description,
		}

		err = c.service.AddCommercialActivity(r.Context(), id, activity)
		if err != nil {
			slog.Error("failed to add commercial activity",
				slog.String("Error", err.Error()),
				slog.String("companyID", id))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _addCommercialActivityError)
			return
		}

		response := CommercialActivityResponse{
			ID:          activity.ID,
			Code:        activity.Code,
			Description: activity.Description,
		}

		httpserver.ReplyJSONResponse(w, http.StatusCreated, response)
	}
}

func (c *CompanyController) removeCommercialActivity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/companies/")
		parts := strings.Split(path, "/commercial-activities/")
		if len(parts) != 2 {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "invalid path format")
			return
		}

		companyID := parts[0]
		activityID := parts[1]

		if companyID == "" || activityID == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "company id and activity id are required")
			return
		}

		err := c.service.RemoveCommercialActivity(r.Context(), companyID, activityID)
		if err != nil {
			slog.Error("failed to remove commercial activity",
				slog.String("Error", err.Error()),
				slog.String("companyID", companyID),
				slog.String("activityID", activityID))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _removeCommercialActivityError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *CompanyController) getCommercialActivities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/companies/")
		path = strings.TrimSuffix(path, "/commercial-activities")
		if path == "" {
			httpserver.ReplyWithError(w, http.StatusBadRequest, "company id is required")
			return
		}

		activities, err := c.service.GetCommercialActivities(r.Context(), path)
		if err != nil {
			slog.Error("failed to get commercial activities",
				slog.String("Error", err.Error()),
				slog.String("companyID", path))
			httpserver.ReplyWithError(w, http.StatusInternalServerError, _getCommercialActivitiesError)
			return
		}

		response := make([]CommercialActivityResponse, len(activities))
		for i, activity := range activities {
			response[i] = CommercialActivityResponse{
				ID:          activity.ID,
				Code:        activity.Code,
				Description: activity.Description,
			}
		}

		httpserver.ReplyJSONResponse(w, http.StatusOK, response)
	}
}

type CompanyCreateRequest struct {
	Name                  string                      `json:"name"`
	Code                  string                      `json:"code"`
	Address               string                      `json:"address"`
	FacturaMovilCompanyID uint64                      `json:"factura_movil_company_id"`
	CommercialActivities  []domain.CommercialActivity `json:"commercial_activities"`
}

type CompanyUpdateRequest struct {
	Name                  string                      `json:"name"`
	Code                  string                      `json:"code"`
	Address               string                      `json:"address"`
	FacturaMovilCompanyID uint64                      `json:"factura_movil_company_id"`
	CommercialActivities  []domain.CommercialActivity `json:"commercial_activities"`
}

type CompanyResponse struct {
	ID                    string                       `json:"id"`
	Name                  string                       `json:"name"`
	Code                  string                       `json:"code"`
	Address               string                       `json:"address"`
	FacturaMovilCompanyID uint64                       `json:"factura_movil_company_id"`
	CommercialActivities  []CommercialActivityResponse `json:"commercial_activities"`
}

type CommercialActivityCreateRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type CommercialActivityResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}
