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
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
)

type PVZService interface {
	Create(ctx context.Context, params application.CreateParams) (*pvz_domain.PVZ, error)
	DeleteLastProduct(ctx context.Context, params application.DeleteLastProductParams) error
	CloseLastReception(ctx context.Context, params application.CloseLastReceptionParams) (*reception_domain.Reception, error)
	ListWithReceptions(ctx context.Context, params application.ListWithReceptionsParams) ([]*pvz_domain.PVZWithReceptions, error)
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
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.JSONError(w, http.StatusForbidden, errors.New("invalid request body"))
		return
	}

	claims, ok := r.Context().Value(httpcommon.DefaultUserKey).(*httpcommon.Claims)
	if !ok {
		httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
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
		switch {
		case errors.Is(err, pvz_domain.ErrAccessDenied):
			httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		case errors.Is(err, pvz_domain.ErrPVZAlreadyExists):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("pvz already exists"))
		default:
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		}
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

	claims, ok := r.Context().Value(httpcommon.DefaultUserKey).(*httpcommon.Claims)
	if !ok {
		httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	params := application.CloseLastReceptionParams{
		PVZID:    pvzID,
		UserRole: auth_domain.Role(claims.Role),
	}

	reception, err := h.svc.CloseLastReception(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, pvz_domain.ErrAccessDenied):
			httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		case errors.Is(err, reception_domain.ErrNoOpenReception):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("reception already closed"))
		default:
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		}
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

	claims, ok := r.Context().Value(httpcommon.DefaultUserKey).(*httpcommon.Claims)
	if !ok {
		httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	params := application.DeleteLastProductParams{
		PVZID:    pvzID,
		UserRole: auth_domain.Role(claims.Role),
	}

	err := h.svc.DeleteLastProduct(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, pvz_domain.ErrAccessDenied):
			httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		case errors.Is(err, reception_domain.ErrNoOpenReception):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("no open reception found"))
		case errors.Is(err, product_domain.ErrNoProductsToDelete):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("no products to delete"))
		default:
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		}
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

	resp := make([]*ListResponse, 0)

	if startDateStr != "" {
		t, err := time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			httpcommon.JSONResponse(w, http.StatusOK, resp)
			return
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse(time.DateOnly, endDateStr)
		if err != nil {
			httpcommon.JSONResponse(w, http.StatusOK, resp)
			return
		}
		endDate = &t
	}

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			httpcommon.JSONResponse(w, http.StatusOK, resp)
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			httpcommon.JSONResponse(w, http.StatusOK, resp)
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
		httpcommon.JSONResponse(w, http.StatusOK, resp)
		return
	}

	// TODO: separate
	for _, val := range result {
		resp = append(resp, &ListResponse{
			PVZ: pvz{
				ID:               *val.PVZ.ID,
				RegistrationDate: *val.PVZ.RegistrationDate,
				City:             string(val.PVZ.City),
			},
			Receptions: func() []receptionWithProducts {
				var receptions []receptionWithProducts
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
					receptions = append(receptions, receptionWithProducts{
						Reception: reception{
							ID:       rec.Reception.ID,
							DateTime: rec.Reception.DateTime,
							PVZID:    rec.Reception.PVZID,
							Status:   string(rec.Reception.Status),
						},
						Products: products,
					})
				}
				return receptions
			}(),
		})
	}

	httpcommon.JSONResponse(w, http.StatusOK, resp)
}
