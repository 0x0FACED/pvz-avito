package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/reception/application"
	reception_http "github.com/0x0FACED/pvz-avito/internal/reception/delivery/http"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/0x0FACED/pvz-avito/internal/reception/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReceptionHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	tests := []struct {
		name           string
		request        reception_http.CreateRequest
		userRole       string
		mockSetup      func(*mocks.MockReceptionService)
		expectedStatus int
		expectError    string
	}{
		{
			name: "successful reception creation",
			request: reception_http.CreateRequest{
				PVZID: "pvz-123",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockReceptionService) {
				m.EXPECT().Create(
					gomock.Any(),
					application.CreateParams{
						PVZID:    "pvz-123",
						UserRole: auth_domain.Role("employee"),
					},
				).Return(&reception_domain.Reception{
					ID:       "rec-123",
					DateTime: now,
					PVZID:    "pvz-123",
					Status:   reception_domain.InProgress,
				}, nil)
			},
			expectedStatus: nethttp.StatusCreated,
		},
		{
			name: "access denied for moderator",
			request: reception_http.CreateRequest{
				PVZID: "pvz-123",
			},
			userRole: "moderator",
			mockSetup: func(m *mocks.MockReceptionService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, reception_domain.ErrAccessDenied)
			},
			expectedStatus: nethttp.StatusForbidden,
			expectError:    "access denied",
		},
		{
			name: "reception already exists",
			request: reception_http.CreateRequest{
				PVZID: "pvz-123",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockReceptionService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, reception_domain.ErrFoundOpenedReception)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectError:    "reception already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receptionSvcMock := mocks.NewMockReceptionService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(receptionSvcMock)
			}

			handler := reception_http.NewHandler(receptionSvcMock)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(nethttp.MethodPost, "/receptions", bytes.NewReader(body))

			ctx := context.WithValue(req.Context(), httpcommon.DefaultUserKey, &httpcommon.Claims{
				Role: tt.userRole,
			})
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectError != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Contains(t, errResp.Error(), tt.expectError)
			} else {
				var resp reception_http.CreateResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.request.PVZID, resp.PVZID)
			}
		})
	}
}

func TestReceptionHandler_Create_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid JSON body", func(t *testing.T) {
		receptionSvcMock := mocks.NewMockReceptionService(ctrl)
		handler := reception_http.NewHandler(receptionSvcMock)

		req := httptest.NewRequest(nethttp.MethodPost, "/receptions", bytes.NewReader([]byte("{invalid}")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusBadRequest, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})

	t.Run("missing claims in context", func(t *testing.T) {
		receptionSvcMock := mocks.NewMockReceptionService(ctrl)
		handler := reception_http.NewHandler(receptionSvcMock)

		validRequest := reception_http.CreateRequest{
			PVZID: "pvz-123",
		}
		body, _ := json.Marshal(validRequest)

		req := httptest.NewRequest(nethttp.MethodPost, "/receptions", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})

	t.Run("invalid claims type in context", func(t *testing.T) {
		receptionSvcMock := mocks.NewMockReceptionService(ctrl)
		handler := reception_http.NewHandler(receptionSvcMock)

		validRequest := reception_http.CreateRequest{
			PVZID: "pvz-123",
		}
		body, _ := json.Marshal(validRequest)

		ctx := context.WithValue(context.Background(), httpcommon.DefaultUserKey, "not-a-claims-object")
		req := httptest.NewRequest(nethttp.MethodPost, "/receptions", bytes.NewReader(body)).WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})
}
