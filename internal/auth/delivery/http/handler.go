package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	"github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
)

type AuthService interface {
	Register(ctx context.Context, params application.RegisterParams) (*domain.User, error)
	Login(ctx context.Context, params application.LoginParams) (*domain.User, error)
}

type Handler struct {
	svc AuthService

	jwtManager *httpcommon.JWTManager
}

func NewHandler(svc AuthService, jwt *httpcommon.JWTManager) *Handler {
	return &Handler{
		svc:        svc,
		jwtManager: jwt,
	}
}

func (h Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /dummyLogin", h.DummyLogin)
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := application.RegisterParams{
		Email:    domain.Email(req.Email),
		Password: req.Password,
		Role:     domain.Role(req.Role),
	}

	if err := params.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	_ = user

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"registered successfully"}`))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := application.LoginParams{
		Email:    domain.Email(req.Email),
		Password: req.Password,
	}

	if err := params.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.Login(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	token, err := h.jwtManager.Generate(user.Email.String(), user.Role.String())
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	type dummyLoginRequest struct {
		Role string `json:"role"`
	}

	var req dummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role != "moderator" && req.Role != "employee" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	token, err := h.jwtManager.GenerateDummy(req.Role)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
