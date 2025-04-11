package application_test

import (
	"context"
	"testing"
	"time"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	product_mocks "github.com/0x0FACED/pvz-avito/internal/product/mocks"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	pvz_mocks "github.com/0x0FACED/pvz-avito/internal/pvz/mocks"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	reception_mocks "github.com/0x0FACED/pvz-avito/internal/reception/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPVZService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regDate := time.Now()
	pvzID := uuid.NewString()
	city := pvz_domain.City("Москва")

	tests := []struct {
		name      string
		params    application.CreateParams
		mockSetup func(*pvz_mocks.MockPVZRepository)
		expectErr error
	}{
		{
			name: "successful creation with provided id and regDate",
			params: application.CreateParams{
				ID:               &pvzID,
				RegistrationDate: &regDate,
				City:             pvz_domain.City(city),
				UserRole:         auth_domain.RoleModerator,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p *pvz_domain.PVZ) (*pvz_domain.PVZ, error) {
						assert.Equal(t, pvzID, *p.ID)
						assert.Equal(t, regDate, *p.RegistrationDate)
						assert.Equal(t, city, p.City)
						return p, nil
					})
			},
			expectErr: nil,
		},
		{
			name: "successful creation with city only",
			params: application.CreateParams{
				City:     pvz_domain.City(city),
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p *pvz_domain.PVZ) (*pvz_domain.PVZ, error) {
						assert.NotNil(t, p.ID)
						assert.NotNil(t, p.RegistrationDate)
						assert.Equal(t, city, p.City)
						return p, nil
					})
			},
			expectErr: nil,
		},
		{
			name: "pvz already exists",
			params: application.CreateParams{
				ID:       &pvzID,
				City:     pvz_domain.City(city),
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, pvz_domain.ErrPVZAlreadyExists)
			},
			expectErr: pvz_domain.ErrPVZAlreadyExists,
		},
		{
			name: "database error",
			params: application.CreateParams{
				City:     pvz_domain.City(city),
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, pvz_domain.ErrInternalDatabase)
			},
			expectErr: pvz_domain.ErrInternalDatabase,
		},
		{
			name: "invalid params - empty city",
			params: application.CreateParams{
				City:     "",
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
			},
			expectErr: pvz_domain.ErrUnsupportedCity,
		},
		{
			name: "invalid params - user not moderator",
			params: application.CreateParams{
				City:     city,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
			},
			expectErr: pvz_domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzRepo := pvz_mocks.NewMockPVZRepository(ctrl)
			receptionRepo := reception_mocks.NewMockReceptionRepository(ctrl)
			productRepo := product_mocks.NewMockProductRepository(ctrl)
			tt.mockSetup(pvzRepo)

			logger := logger.NewTestLogger()

			service := application.NewPVZService(pvzRepo, receptionRepo, productRepo, logger)
			_, err := service.Create(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPVZService_CloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()
	receptionID := uuid.NewString()

	tests := []struct {
		name      string
		params    application.CloseLastReceptionParams
		mockSetup func(*reception_mocks.MockReceptionRepository)
		expectErr error
	}{
		{
			name: "successful close",
			params: application.CloseLastReceptionParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					CloseLastReception(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     receptionID,
						PVZID:  pvzID,
						Status: reception_domain.Close,
					}, nil)
			},
			expectErr: nil,
		},
		{
			name: "no open reception",
			params: application.CloseLastReceptionParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					CloseLastReception(gomock.Any(), pvzID).
					Return(nil, reception_domain.ErrNoOpenReception)
			},
			expectErr: reception_domain.ErrNoOpenReception,
		},
		{
			name: "database error",
			params: application.CloseLastReceptionParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					CloseLastReception(gomock.Any(), pvzID).
					Return(nil, pvz_domain.ErrInternalDatabase)
			},
			expectErr: pvz_domain.ErrInternalDatabase,
		},
		{
			name: "invalid params - empty pvzID",
			params: application.CloseLastReceptionParams{
				PVZID:    "",
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
			},
			expectErr: pvz_domain.ErrInvalidIDFormat,
		},
		{
			name: "invalid params - not employee",
			params: application.CloseLastReceptionParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
			},
			expectErr: pvz_domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzRepo := pvz_mocks.NewMockPVZRepository(ctrl)
			receptionRepo := reception_mocks.NewMockReceptionRepository(ctrl)
			productRepo := product_mocks.NewMockProductRepository(ctrl)
			tt.mockSetup(receptionRepo)

			logger := logger.NewTestLogger()

			service := application.NewPVZService(pvzRepo, receptionRepo, productRepo, logger)
			_, err := service.CloseLastReception(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPVZService_DeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()
	receptionID := uuid.NewString()

	tests := []struct {
		name      string
		params    application.DeleteLastProductParams
		mockSetup func(*reception_mocks.MockReceptionRepository, *product_mocks.MockProductRepository)
		expectErr error
	}{
		{
			name: "successful delete",
			params: application.DeleteLastProductParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     receptionID,
						PVZID:  pvzID,
						Status: reception_domain.InProgress,
					}, nil)

				p.EXPECT().
					DeleteLastFromReception(gomock.Any(), receptionID).
					Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "no open reception",
			params: application.DeleteLastProductParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, reception_domain.ErrNoOpenReception)
			},
			expectErr: reception_domain.ErrNoOpenReception,
		},
		{
			name: "error finding reception",
			params: application.DeleteLastProductParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, pvz_domain.ErrInternalDatabase)
			},
			expectErr: pvz_domain.ErrInternalDatabase,
		},
		{
			name: "database error",
			params: application.DeleteLastProductParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     receptionID,
						PVZID:  pvzID,
						Status: reception_domain.InProgress,
					}, nil)

				p.EXPECT().
					DeleteLastFromReception(gomock.Any(), receptionID).
					Return(pvz_domain.ErrInternalDatabase)
			},
			expectErr: pvz_domain.ErrInternalDatabase,
		},
		{
			name: "invalid params - empty pvzID",
			params: application.DeleteLastProductParams{
				PVZID:    "",
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
			},
			expectErr: pvz_domain.ErrInvalidIDFormat,
		},
		{
			name: "invalid params - not employee",
			params: application.DeleteLastProductParams{
				PVZID:    pvzID,
				UserRole: auth_domain.RoleModerator,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
			},
			expectErr: pvz_domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzRepo := pvz_mocks.NewMockPVZRepository(ctrl)
			receptionRepo := reception_mocks.NewMockReceptionRepository(ctrl)
			productRepo := product_mocks.NewMockProductRepository(ctrl)
			tt.mockSetup(receptionRepo, productRepo)

			log := logger.NewTestLogger()
			service := application.NewPVZService(pvzRepo, receptionRepo, productRepo, log)

			err := service.DeleteLastProduct(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPVZService_ListWithReceptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	receptionID := uuid.NewString()
	pvzID := uuid.NewString()
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	page := 1
	limit := 10

	tests := []struct {
		name        string
		params      application.ListWithReceptionsParams
		mockSetup   func(*pvz_mocks.MockPVZRepository)
		expectCount int
		expectErr   error
	}{
		{
			name: "successful list",
			params: application.ListWithReceptionsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      &page,
				Limit:     &limit,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					ListWithReceptions(gomock.Any(), &startDate, &endDate, page, limit).
					Return([]*pvz_domain.PVZWithReceptions{
						{
							PVZ: &pvz_domain.PVZ{
								ID:               &pvzID,
								RegistrationDate: &startDate,
								City:             "Moscow",
							},
							Receptions: []*pvz_domain.ReceptionWithProducts{
								{
									Reception: &reception_domain.Reception{
										ID:     receptionID,
										PVZID:  pvzID,
										Status: reception_domain.InProgress,
									},
									Products: []*product_domain.Product{
										{
											ID:          uuid.NewString(),
											DateTime:    time.Now().Add(-1 * time.Hour),
											Type:        product_domain.Clothes,
											ReceptionID: receptionID,
										},
									},
								},
							},
						},
					}, nil)
			},
			expectCount: 1,
			expectErr:   nil,
		},
		{
			name: "empty result",
			params: application.ListWithReceptionsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      &page,
				Limit:     &limit,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					ListWithReceptions(gomock.Any(), &startDate, &endDate, page, limit).
					Return([]*pvz_domain.PVZWithReceptions{}, nil)
			},
			expectCount: 0,
			expectErr:   nil,
		},
		{
			name: "database error",
			params: application.ListWithReceptionsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      &page,
				Limit:     &limit,
			},
			mockSetup: func(r *pvz_mocks.MockPVZRepository) {
				r.EXPECT().
					ListWithReceptions(gomock.Any(), &startDate, &endDate, page, limit).
					Return(nil, pvz_domain.ErrInternalDatabase)
			},
			expectCount: 0,
			expectErr:   pvz_domain.ErrInternalDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvzRepo := pvz_mocks.NewMockPVZRepository(ctrl)
			receptionRepo := reception_mocks.NewMockReceptionRepository(ctrl)
			productRepo := product_mocks.NewMockProductRepository(ctrl)
			tt.mockSetup(pvzRepo)

			log := logger.NewTestLogger()
			service := application.NewPVZService(pvzRepo, receptionRepo, productRepo, log)

			result, err := service.ListWithReceptions(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectCount, len(result))
			}
		})
	}
}
