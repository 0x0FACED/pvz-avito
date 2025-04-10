package http

import "net/http"

type ReceptionService interface {
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
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"not impl"}`))
}
