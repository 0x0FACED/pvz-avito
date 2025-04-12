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
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_http "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/http"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	"github.com/0x0FACED/pvz-avito/internal/pvz/mocks"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPVZHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()
	now := time.Now()
	tests := []struct {
		name           string
		request        pvz_http.CreateRequest
		userRole       string
		mockSetup      func(*mocks.MockPVZService)
		expectedStatus int
		expectErr      string
	}{
		{
			name: "successful creation",
			request: pvz_http.CreateRequest{
				ID:               &pvzID,
				RegistrationDate: &now,
				City:             "Москва",
			},
			userRole: "moderator",
			mockSetup: func(m *mocks.MockPVZService) {
				expectedParams := application.CreateParams{
					ID:               &pvzID,
					RegistrationDate: &now,
					City:             pvz_domain.City("Москва"),
					UserRole:         auth_domain.Role("moderator"),
				}
				m.EXPECT().Create(
					gomock.Any(),
					createParamsMatcher{expectedParams},
				).Return(&pvz_domain.PVZ{
					ID:               &pvzID,
					RegistrationDate: &now,
					City:             pvz_domain.City("Москва"),
				}, nil)
			},
			expectedStatus: nethttp.StatusCreated,
		},
		{
			name: "access denied for employee",
			request: pvz_http.CreateRequest{
				ID:               &pvzID,
				RegistrationDate: &now,
				City:             "Москва",
			},
			userRole: "employee",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, pvz_domain.ErrAccessDenied)
			},
			expectedStatus: nethttp.StatusForbidden,
			expectErr:      "access denied",
		},
		{
			name: "pvz already exists",
			request: pvz_http.CreateRequest{
				ID:               &pvzID,
				RegistrationDate: &now,
				City:             "Москва",
			},
			userRole: "moderator",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, pvz_domain.ErrPVZAlreadyExists)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzSvcMock := mocks.NewMockPVZService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(pvzSvcMock)
			}

			handler := pvz_http.NewHandler(pvzSvcMock)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(nethttp.MethodPost, "/pvz", bytes.NewReader(body))

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
				var resp pvz_http.CreateResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.Equal(t, *tt.request.ID, *resp.ID)
				assert.Equal(t, tt.request.City, resp.City)
			}
		})
	}

	t.Run("invalid JSON body", func(t *testing.T) {
		pvzSvcMock := mocks.NewMockPVZService(ctrl)
		handler := pvz_http.NewHandler(pvzSvcMock)

		req := httptest.NewRequest(nethttp.MethodPost, "/pvz", bytes.NewReader([]byte("{invalid}")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "invalid request body", errResp.Error())
	})

	t.Run("missing claims in context", func(t *testing.T) {
		pvzSvcMock := mocks.NewMockPVZService(ctrl)
		handler := pvz_http.NewHandler(pvzSvcMock)

		validRequest := pvz_http.CreateRequest{
			ID:               &pvzID,
			RegistrationDate: &now,
			City:             "Москва",
		}
		body, _ := json.Marshal(validRequest)

		req := httptest.NewRequest(nethttp.MethodPost, "/pvz", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})
}

func TestPVZHandler_CloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	tests := []struct {
		name           string
		pvzID          string
		userRole       string
		mockSetup      func(*mocks.MockPVZService)
		expectedStatus int
		expectErr      string
	}{
		{
			name:     "successful close",
			pvzID:    "pvz-123",
			userRole: "employee",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().CloseLastReception(
					gomock.Any(),
					application.CloseLastReceptionParams{
						PVZID:    "pvz-123",
						UserRole: auth_domain.RoleEmployee,
					},
				).Return(&reception_domain.Reception{
					ID:       "rec-123",
					DateTime: now,
					PVZID:    "pvz-123",
					Status:   reception_domain.Close,
				}, nil)
			},
			expectedStatus: nethttp.StatusOK,
		},
		{
			name:     "access denied for moderator",
			pvzID:    "pvz-123",
			userRole: "moderator",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).
					Return(nil, pvz_domain.ErrAccessDenied)
			},
			expectedStatus: nethttp.StatusForbidden,
			expectErr:      "access denied",
		},
		{
			name:     "reception already closed",
			pvzID:    "pvz-123",
			userRole: "employee",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).
					Return(nil, reception_domain.ErrNoOpenReception)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "reception already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzSvcMock := mocks.NewMockPVZService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(pvzSvcMock)
			}

			handler := pvz_http.NewHandler(pvzSvcMock)

			req := httptest.NewRequest("POST", "/pvz/"+tt.pvzID+"/close_last_reception", nil)

			req.SetPathValue("pvzId", "pvz-123")
			ctx := context.WithValue(req.Context(), httpcommon.DefaultUserKey, &httpcommon.Claims{
				Role: tt.userRole,
			})
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			handler.CloseLastReception(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Equal(t, tt.expectErr, errResp.Error())
			} else {
				var resp pvz_http.CloseResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.Equal(t, tt.pvzID, resp.PVZID)
				assert.Equal(t, reception_domain.Close, reception_domain.Status(resp.Status))
			}
		})
	}

	t.Run("missing claims in context", func(t *testing.T) {
		pvzSvcMock := mocks.NewMockPVZService(ctrl)
		handler := pvz_http.NewHandler(pvzSvcMock)

		req := httptest.NewRequest(nethttp.MethodPost, "/pvz/pvz-123/close_last_reception", nil)
		rec := httptest.NewRecorder()

		handler.CloseLastReception(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})
}

func TestPVZHandler_DeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		pvzID          string
		userRole       string
		mockSetup      func(*mocks.MockPVZService)
		expectedStatus int
		expectErr      string
	}{
		{
			name:     "successful delete",
			pvzID:    "pvz-123",
			userRole: "employee",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().DeleteLastProduct(
					gomock.Any(),
					application.DeleteLastProductParams{
						PVZID:    "pvz-123",
						UserRole: auth_domain.RoleEmployee,
					},
				).Return(nil)
			},
			expectedStatus: nethttp.StatusOK,
		},
		{
			name:     "access denied for moderator",
			pvzID:    "pvz-123",
			userRole: "moderator",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Return(pvz_domain.ErrAccessDenied)
			},
			expectedStatus: nethttp.StatusForbidden,
			expectErr:      "access denied",
		},
		{
			name:     "no open reception",
			pvzID:    "pvz-123",
			userRole: "employee",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Return(reception_domain.ErrNoOpenReception)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "no open reception found",
		},
		{
			name:     "no products to delete",
			pvzID:    "pvz-123",
			userRole: "moderator",
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Return(product_domain.ErrNoProductsToDelete)
			},
			expectedStatus: nethttp.StatusBadRequest,
			expectErr:      "no products to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzSvcMock := mocks.NewMockPVZService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(pvzSvcMock)
			}

			handler := pvz_http.NewHandler(pvzSvcMock)

			req := httptest.NewRequest("POST", "/pvz/"+tt.pvzID+"/delete_last_product", nil)

			req.SetPathValue("pvzId", "pvz-123")
			ctx := context.WithValue(req.Context(), httpcommon.DefaultUserKey, &httpcommon.Claims{
				Role: tt.userRole,
			})
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			handler.DeleteLastProduct(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectErr != "" {
				var errResp httpcommon.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&errResp)
				assert.Equal(t, tt.expectErr, errResp.Error())
			}
		})
	}

	t.Run("missing claims in context", func(t *testing.T) {
		pvzSvcMock := mocks.NewMockPVZService(ctrl)
		handler := pvz_http.NewHandler(pvzSvcMock)

		req := httptest.NewRequest(nethttp.MethodPost, "/pvz/pvz-123/delete_last_product", nil)
		rec := httptest.NewRecorder()

		handler.DeleteLastProduct(rec, req)

		assert.Equal(t, nethttp.StatusForbidden, rec.Code)
		var errResp httpcommon.ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&errResp)
		assert.Equal(t, "access denied", errResp.Error())
	})
}

func TestPVZHandler_ListWithReceptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()

	now := time.Date(2025, 4, 12, 0, 0, 0, 0, time.UTC)
	startDate := now.AddDate(0, -1, 0)
	endDate := now
	page := 1
	limit := 10

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(*mocks.MockPVZService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful list with all params",
			queryParams: map[string]string{
				"startDate": startDate.Format("2006-01-02"),
				"endDate":   endDate.Format("2006-01-02"),
				"page":      "1",
				"limit":     "10",
			},
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().ListWithReceptions(
					gomock.Any(),
					application.ListWithReceptionsParams{
						StartDate: &startDate,
						EndDate:   &endDate,
						Page:      &page,
						Limit:     &limit,
					},
				).Return([]*pvz_domain.PVZWithReceptions{
					{
						PVZ: &pvz_domain.PVZ{
							ID:               &pvzID,
							RegistrationDate: &now,
							City:             "Москва",
						},
						Receptions: []*pvz_domain.ReceptionWithProducts{
							{
								Reception: &reception_domain.Reception{
									ID:       "rec-1",
									DateTime: now,
									PVZID:    "pvz-1",
									Status:   reception_domain.InProgress,
								},
								Products: []*product_domain.Product{
									{
										ID:          "prod-1",
										DateTime:    now,
										Type:        product_domain.Electronics,
										ReceptionID: "rec-1",
									},
								},
							},
						},
					},
				}, nil)
			},
			expectedStatus: nethttp.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "successful list with no params",
			queryParams: map[string]string{},
			mockSetup: func(m *mocks.MockPVZService) {
				m.EXPECT().ListWithReceptions(
					gomock.Any(),
					application.ListWithReceptionsParams{
						Page:  &page,
						Limit: &limit,
					},
				).Return([]*pvz_domain.PVZWithReceptions{}, nil)
			},
			expectedStatus: nethttp.StatusOK,
			expectedCount:  0,
		},
		{
			name: "invalid date format",
			queryParams: map[string]string{
				"startDate": "invalid-date",
			},
			mockSetup:      nil,
			expectedStatus: nethttp.StatusOK,
			expectedCount:  0,
		},
		{
			name: "invalid page number",
			queryParams: map[string]string{
				"page": "invalid",
			},
			mockSetup:      nil,
			expectedStatus: nethttp.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzSvcMock := mocks.NewMockPVZService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(pvzSvcMock)
			}

			handler := pvz_http.NewHandler(pvzSvcMock)

			req := httptest.NewRequest(nethttp.MethodGet, "/pvz", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			handler.ListWithReceptions(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var resp []pvz_http.ListResponse
			_ = json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tt.expectedCount, len(resp))
		})
	}
}

type createParamsMatcher struct {
	expected application.CreateParams
}

func (m createParamsMatcher) Matches(x interface{}) bool {
	actual, ok := x.(application.CreateParams)
	if !ok {
		return false
	}

	if actual.ID == nil || m.expected.ID == nil || *actual.ID != *m.expected.ID {
		return false
	}

	if actual.RegistrationDate == nil || m.expected.RegistrationDate == nil ||
		!actual.RegistrationDate.Equal(*m.expected.RegistrationDate) {
		return false
	}

	if actual.City != m.expected.City || actual.UserRole != m.expected.UserRole {
		return false
	}

	return true
}

func (m createParamsMatcher) String() string {
	return "matches application.CreateParams by value"
}
