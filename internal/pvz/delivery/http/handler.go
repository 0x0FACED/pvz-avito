package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
)

type PVZService interface {
	Create(ctx context.Context, params application.CreateParams) (*pvz_domain.PVZ, error)
	FindByID(ctx context.Context, id string) (*pvz_domain.PVZ, error)
	DeleteLastProduct(ctx context.Context, params application.DeleteLastProductParams) error
	ListWithReceptions(ctx context.Context, params application.ListWithReceptionsParams) ([]*pvz_domain.PVZWithReceptions, error)
	CloseLastReception(ctx context.Context, params application.CloseLastReceptionParams) (*reception_domain.Reception, error)
}

type Handler struct {
	svc PVZService
}

func NewHandler(svc PVZService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /pvz", h.Create)
	mux.HandleFunc("GET /pvz", h.ListWithReceptions)
	mux.HandleFunc("POST /pvz/{pvzId}/close_last_reception", h.CloseLastReception)
	mux.HandleFunc("POST /pvz/{pvzId}/delete_last_product", h.DeleteLastProduct)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("user").(*httpcommon.Claims)
	if !ok {
		http.Error(w, "User not found in context", http.StatusBadRequest)
		return
	}

	params := application.CreateParams{
		ID:               req.ID,
		RegistrationDate: req.RegistrationDate,
		City:             pvz_domain.City(req.City),
		UserRole:         auth_domain.Role(claims.Role),
	}

	pvz, err := h.svc.Create(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	resp := CreateResponse{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City.String(),
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}

func (h *Handler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	pvzID := r.PathValue("pvzId")

	claims, ok := r.Context().Value("user").(*httpcommon.Claims)
	if !ok {
		http.Error(w, "User not found in context", http.StatusBadRequest)
		return
	}

	params := application.CloseLastReceptionParams{
		PVZID:    pvzID,
		UserRole: auth_domain.Role(claims.Role),
	}

	reception, err := h.svc.CloseLastReception(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	resp := CloseResponse{
		ID:       reception.ID,
		DateTime: reception.DateTime,
		PVZID:    reception.PVZID,
		Status:   reception.Status.String(),
	}

	httpcommon.JSONResponse(w, http.StatusOK, resp)
}

func (h *Handler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	pvzID := r.PathValue("pvzId")

	claims, ok := r.Context().Value("user").(*httpcommon.Claims)
	if !ok {
		http.Error(w, "User not found in context", http.StatusBadRequest)
		return
	}

	params := application.DeleteLastProductParams{
		PVZID:    pvzID,
		UserRole: auth_domain.Role(claims.Role),
	}

	err := h.svc.DeleteLastProduct(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	httpcommon.EmptyResponse(w, http.StatusOK)
}

func (h *Handler) ListWithReceptions(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var (
		startDate *time.Time
		endDate   *time.Time
		page      = 1
		limit     = 10
	)

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("start date must be YYYY-MM-DD format"))
			return
		}
		startDate = &t
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			// TODO: handle err
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("end date must be YYYY-MM-DD format"))
			return
		}
		endDate = &t
	}

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			// TODO: handle err
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid page value"))
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			// TODO: handle err
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid limit value"))
			return
		}
		limit = l
	}

	params := application.ListWithReceptionsParams{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      &page,
		Limit:     &limit,
	}

	result, err := h.svc.ListWithReceptions(r.Context(), params)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]*ListResponse, 0, len(result))

	// TODO: separate
	for _, val := range result {
		resp = append(resp, &ListResponse{
			PVZ: pvz{
				ID:               *val.PVZ.ID,
				RegistrationDate: *val.PVZ.RegistrationDate,
				City:             string(val.PVZ.City),
			},
			Receptions: func() []reception {
				var receptions []reception
				for _, rec := range val.Receptions {
					var products []product
					for _, prod := range rec.Products {
						products = append(products, product{
							ID:          prod.ID,
							DateTime:    prod.DateTime,
							Type:        string(prod.Type),
							ReceptionID: prod.ReceptionID,
						})
					}
					receptions = append(receptions, reception{
						ID:       rec.Reception.ID,
						DateTime: rec.Reception.DateTime,
						PVZID:    rec.Reception.PVZID,
						Status:   string(rec.Reception.Status),
						Products: products,
					})
				}
				return receptions
			}(),
		})
	}

	httpcommon.JSONResponse(w, http.StatusOK, resp)
}
