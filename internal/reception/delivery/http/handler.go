package http

import (
	"context"
	"encoding/json"
	"net/http"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/reception/application"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
)

type ReceptionService interface {
	Create(ctx context.Context, params application.CreateParams) (*reception_domain.Reception, error)
}

type Handler struct {
	svc ReceptionService
}

func NewHandler(svc ReceptionService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /receptions", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	type createRequest struct {
		PVZID string `json:"pvzId"`
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
		PVZID:    req.PVZID,
		UserRole: auth_domain.Role(claims.Role),
	}

	reception, err := h.svc.Create(r.Context(), params)
	if err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, err)
		return
	}

	httpcommon.JSONResponse(w, http.StatusCreated, reception)
}
