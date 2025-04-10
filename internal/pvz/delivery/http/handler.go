package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	"github.com/0x0FACED/pvz-avito/internal/pvz/domain"
)

type PVZService interface {
	Create(ctx context.Context, params application.CreateParams) (*domain.PVZ, error)
	FindByID(ctx context.Context, id string) (*domain.PVZ, error)
	ListWithReceptions(ctx context.Context, params application.ListWithReceptionsParams) ([]*domain.PVZWithReceptions, error)
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
	type createRequest struct {
		ID               *string    `json:"id,omitempty"`
		RegistrationDate *time.Time `json:"registrationDate,omitempty"`
		City             string     `json:"city"`
	}

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
		City:             domain.City(req.City),
		UserRole:         claims.Role,
	}

	pvz, err := h.svc.Create(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	type createResponse struct {
		ID               string      `json:"id"`
		RegistrationDate time.Time   `json:"registrationDate"`
		City             domain.City `json:"city"`
	}

	resp := createResponse{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}

func (h *Handler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"not impl"}`))
}

func (h *Handler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"not impl"}`))
}

func (h *Handler) ListWithReceptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"not impl"}`))
}
