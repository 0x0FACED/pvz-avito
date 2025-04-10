package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
	type createRequest struct {
		Type  string `json:"type"`
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
		Type:     product_domain.ProductType(req.Type),
		PVZID:    req.PVZID,
		UserRole: auth_domain.Role(claims.Role),
	}

	product, err := h.svc.Create(r.Context(), params)
	if err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, err)
		return
	}

	type createResponse struct {
		ID          string                     `json:"id"`
		DateTime    time.Time                  `json:"dateTime"`
		Type        product_domain.ProductType `json:"type"`
		ReceptionID string                     `json:"receptionId"`
	}

	resp := createResponse{
		ID:          product.ID,
		DateTime:    product.DateTime,
		Type:        product.Type,
		ReceptionID: product.ReceptionID,
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}
