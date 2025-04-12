package http_test

import (
	"bytes"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	auth_http "github.com/0x0FACED/pvz-avito/internal/auth/delivery/http"
	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/auth/mocks"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		request        auth_http.RegisterRequest
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		expectErr      string
	}{
		{
			name: "successful registration",
			request: auth_http.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "moderator",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Register(
					gomock.Any(),
					application.RegisterParams{
						Email:    auth_domain.Email("test@example.com"),
						Password: "password123",
						Role:     auth_domain.Role("moderator"),
					},
				).Return(&auth_domain.User{
					ID:    "123",
					Email: auth_domain.Email("test@example.com"),
					Role:  auth_domain.Role("moderator"),
				}, nil)
			},
			expectedStatus: nethttp.StatusCreated,
		},
		{
			name: "invalid email format",
			request: auth_http.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Role:     "moderator",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrInvalidEmail)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "invalid login or password",
		},
		{
			name: "user already exists",
			request: auth_http.RegisterRequest{
				Email:    "exists@example.com",
				Password: "password123",
				Role:     "employee",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrUserAlreadyExists)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      auth_domain.ErrUserAlreadyExists.Error(),
		},
		{
			name: "invalid role",
			request: auth_http.RegisterRequest{
				Email:    "user@example.com",
				Password: "password123",
				Role:     "notavalidrole",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrInvalidRole)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvcMock := mocks.NewMockAuthService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(authSvcMock)
			}

			jwtManager := httpcommon.NewManager("test-secret", 3600)
			handler := auth_http.NewHandler(authSvcMock, jwtManager)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Contains(t, tt.expectErr, errResp.Error())
			} else {
				var resp auth_http.RegisterResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.request.Email, resp.Email)
				assert.Equal(t, tt.request.Role, resp.Role)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		request        auth_http.LoginRequest
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		expectToken    bool
		expectErr      string
	}{
		{
			name: "successful login",
			request: auth_http.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Login(
					gomock.Any(),
					application.LoginParams{
						Email:    auth_domain.Email("test@example.com"),
						Password: "password123",
					},
				).Return(&auth_domain.User{
					ID:       uuid.NewString(),
					Email:    "test@example.com",
					Password: "hashofpassword",
					Role:     "moderator",
				}, nil)
			},
			expectedStatus: nethttp.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid creds",
			request: auth_http.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.EXPECT().Login(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrInvalidPassword)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "invalid login or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvcMock := mocks.NewMockAuthService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(authSvcMock)
			}

			jwtManager := httpcommon.NewManager("test-secret", 3600)
			handler := auth_http.NewHandler(authSvcMock, jwtManager)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Contains(t, tt.expectErr, errResp.Error())
			} else if tt.expectToken {
				token := rec.Body.String()
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthHandler_DummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		request        auth_http.DummyLoginRequest
		expectedStatus int
		expectToken    bool
		expectErr      string
	}{
		{
			name: "successful dummy login with user role",
			request: auth_http.DummyLoginRequest{
				Role: "moderator",
			},
			expectedStatus: nethttp.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid role",
			request: auth_http.DummyLoginRequest{
				Role: "notavalidrole",
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvcMock := mocks.NewMockAuthService(ctrl)
			jwtManager := httpcommon.NewManager("test-secret", 3600)
			handler := auth_http.NewHandler(authSvcMock, jwtManager)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.DummyLogin(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Contains(t, tt.expectErr, errResp.Error())
			} else if tt.expectToken {
				token := rec.Body.String()
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthHandler_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid request body (register)", func(t *testing.T) {
		authSvcMock := mocks.NewMockAuthService(ctrl)
		jwtManager := httpcommon.NewManager("test-secret", 3600)
		handler := auth_http.NewHandler(authSvcMock, jwtManager)

		body, _ := json.Marshal([]byte("{invalid}"))

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, nethttp.StatusBadRequest, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})

	t.Run("invalid request body (login)", func(t *testing.T) {
		authSvcMock := mocks.NewMockAuthService(ctrl)
		jwtManager := httpcommon.NewManager("test-secret", 3600)
		handler := auth_http.NewHandler(authSvcMock, jwtManager)

		body, _ := json.Marshal([]byte("{invalid}"))

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, nethttp.StatusBadRequest, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})

	t.Run("invalid request body (dummyLogin)", func(t *testing.T) {
		authSvcMock := mocks.NewMockAuthService(ctrl)
		jwtManager := httpcommon.NewManager("test-secret", 3600)
		handler := auth_http.NewHandler(authSvcMock, jwtManager)

		body, _ := json.Marshal([]byte("{invalid}"))

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.DummyLogin(rec, req)

		assert.Equal(t, nethttp.StatusBadRequest, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})
}
