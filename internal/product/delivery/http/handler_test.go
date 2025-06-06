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
	"github.com/0x0FACED/pvz-avito/internal/product/application"
	product_http "github.com/0x0FACED/pvz-avito/internal/product/delivery/http"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/0x0FACED/pvz-avito/internal/product/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	tests := []struct {
		name           string
		request        product_http.CreateRequest
		userRole       string
		mockSetup      func(*mocks.MockProductService)
		expectedStatus int
		expectErr      string
	}{
		{
			name: "successful product creation",
			request: product_http.CreateRequest{
				Type:  "электроника",
				PVZID: "pvz-123",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockProductService) {
				m.EXPECT().Create(
					gomock.Any(),
					application.CreateParams{
						Type:     product_domain.ProductType("электроника"),
						PVZID:    "pvz-123",
						UserRole: auth_domain.Role("employee"),
					},
				).Return(&product_domain.Product{
					ID:          "prod-123",
					DateTime:    now,
					Type:        product_domain.ProductType("электроника"),
					ReceptionID: "rec-123",
				}, nil)
			},
			expectedStatus: nethttp.StatusCreated,
		},
		{
			name: "access denied for moderator",
			request: product_http.CreateRequest{
				Type:  "электроника",
				PVZID: "pvz-123",
			},
			userRole: "moderator",
			mockSetup: func(m *mocks.MockProductService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, product_domain.ErrAccessDenied)
			},
			expectedStatus: nethttp.StatusForbidden,
			expectErr:      "access denied",
		},
		{
			name: "invalid product type",
			request: product_http.CreateRequest{
				Type:  "invalidtype",
				PVZID: "pvz-123",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockProductService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, product_domain.ErrInvalidProductType)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "invalid request",
		},
		{
			name: "reception not found",
			request: product_http.CreateRequest{
				Type:  "электроника",
				PVZID: "pvz-123",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockProductService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, product_domain.ErrReceptionNotFound)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "reception not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productSvcMock := mocks.NewMockProductService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(productSvcMock)
			}

			handler := product_http.NewHandler(productSvcMock)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(nethttp.MethodPost, "/products", bytes.NewReader(body))

			ctx := context.WithValue(req.Context(), httpcommon.DefaultUserKey, &httpcommon.Claims{
				Role: tt.userRole,
			})
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Contains(t, errResp.Error(), tt.expectErr)
			} else {
				var resp product_http.CreateResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.request.Type, resp.Type)
			}
		})
	}
}

func TestProductHandler_Create_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid JSON body", func(t *testing.T) {
		productSvcMock := mocks.NewMockProductService(ctrl)
		handler := product_http.NewHandler(productSvcMock)

		req := httptest.NewRequest(nethttp.MethodPost, "/products", bytes.NewReader([]byte("{invalid}")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusBadRequest, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})

	t.Run("missing claims in context", func(t *testing.T) {
		productSvcMock := mocks.NewMockProductService(ctrl)
		handler := product_http.NewHandler(productSvcMock)

		validRequest := product_http.CreateRequest{
			Type:  "электроника",
			PVZID: "pvz-123",
		}
		body, _ := json.Marshal(validRequest)

		req := httptest.NewRequest(nethttp.MethodPost, "/products", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})

	t.Run("invalid claims type in context", func(t *testing.T) {
		productSvcMock := mocks.NewMockProductService(ctrl)
		handler := product_http.NewHandler(productSvcMock)

		validRequest := product_http.CreateRequest{
			Type:  "электроника",
			PVZID: "pvz-123",
		}
		body, _ := json.Marshal(validRequest)

		ctx := context.WithValue(context.Background(), httpcommon.DefaultUserKey, "not-a-claims-object")
		req := httptest.NewRequest(nethttp.MethodPost, "/products", bytes.NewReader(body)).WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})
}
