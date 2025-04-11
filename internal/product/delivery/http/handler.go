package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/product/application"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
)

type ProductService interface {
	Create(ctx context.Context, product application.CreateParams) (*product_domain.Product, error)
}

type Handler struct {
	svc ProductService
}

func NewHandler(svc ProductService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	claims, ok := r.Context().Value("user").(*httpcommon.Claims)
	if !ok {
		httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	params := application.CreateParams{
		Type:     product_domain.ProductType(req.Type),
		PVZID:    req.PVZID,
		UserRole: auth_domain.Role(claims.Role),
	}

	product, err := h.svc.Create(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, product_domain.ErrAccessDenied):
			httpcommon.JSONError(w, http.StatusForbidden, errors.New("access denied"))
		case errors.Is(err, product_domain.ErrReceptionNotFound):
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("reception not found"))
		default:
			httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		}
		return
	}

	resp := CreateResponse{
		ID:          product.ID,
		DateTime:    product.DateTime,
		Type:        product.Type.String(),
		ReceptionID: product.ReceptionID,
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}
