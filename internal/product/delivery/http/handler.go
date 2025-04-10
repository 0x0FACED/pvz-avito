package http

import (
	"context"
	"net/http"

	"github.com/0x0FACED/pvz-avito/internal/product/application"
	"github.com/0x0FACED/pvz-avito/internal/product/domain"
)

type ProductService interface {
	Create(ctx context.Context, product application.CreateParams) (*domain.Product, error)
}

type Handler struct {
	svc ProductService
}

func NewHandler(svc ProductService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"not impl"}`))
}
