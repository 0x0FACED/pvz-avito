package http

import (
	"context"
	"encoding/json"
	"errors"
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

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(httpcommon.DefaultUserKey).(*httpcommon.Claims)
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
		switch {
		case errors.Is(err, reception_domain.ErrAccessDenied):
			httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		case errors.Is(err, reception_domain.ErrFoundOpenedReception):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("reception already exists"))
		default:
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		}
		return
	}

	resp := CreateResponse{
		ID:       reception.ID,
		DateTime: reception.DateTime,
		PVZID:    reception.PVZID,
		Status:   reception.Status.String(),
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}
