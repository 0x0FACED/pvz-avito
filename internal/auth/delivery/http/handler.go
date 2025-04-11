package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
)

type AuthService interface {
	Register(ctx context.Context, params application.RegisterParams) (*auth_domain.User, error)
	Login(ctx context.Context, params application.LoginParams) (*auth_domain.User, error)
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
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := application.RegisterParams{
		Email:    auth_domain.Email(req.Email),
		Password: req.Password,
		Role:     auth_domain.Role(req.Role),
	}

	user, err := h.svc.Register(r.Context(), params)
	if err != nil {
		// TODO: change
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	resp := RegisterResponse{
		ID:    user.ID,
		Email: user.Email.String(),
		Role:  user.Role.String(),
	}

	httpcommon.JSONResponse(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := application.LoginParams{
		Email:    auth_domain.Email(req.Email),
		Password: req.Password,
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

	httpcommon.DefaultResponse(w, http.StatusOK, []byte(token))
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	role := auth_domain.Role(req.Role)

	if err := role.Validate(); err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.jwtManager.GenerateDummy(req.Role)
	if err != nil {
		httpcommon.JSONError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	httpcommon.DefaultResponse(w, http.StatusOK, []byte(token))
}
