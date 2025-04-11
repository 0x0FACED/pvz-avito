package application_test

import (
	"context"
	"testing"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/product/application"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	product_mocks "github.com/0x0FACED/pvz-avito/internal/product/mocks"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	reception_mocks "github.com/0x0FACED/pvz-avito/internal/reception/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProductCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()
	receptionID := uuid.NewString()

	validParams := application.CreateParams{
		PVZID:    pvzID,
		Type:     product_domain.Electronics,
		UserRole: auth_domain.RoleEmployee,
	}

	invalidRoleParams := application.CreateParams{
		PVZID:    pvzID,
		Type:     product_domain.Electronics,
		UserRole: auth_domain.RoleModerator,
	}

	tests := []struct {
		name      string
		params    application.CreateParams
		mockSetup func(*reception_mocks.MockReceptionRepository, *product_mocks.MockProductRepository)
		expectErr error
	}{
		{
			name:   "successful creation",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     receptionID,
						PVZID:  pvzID,
						Status: reception_domain.InProgress,
					}, nil)

				p.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, product *product_domain.Product) (*product_domain.Product, error) {
						assert.Equal(t, receptionID, product.ReceptionID)
						assert.Equal(t, product_domain.Electronics, product.Type)
						return product, nil
					})
			},
			expectErr: nil,
		},
		{
			name:   "no open reception found",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, reception_domain.ErrNoOpenReception)
			},
			expectErr: reception_domain.ErrNoOpenReception,
		},
		{
			name:   "database error when finding reception",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, assert.AnError)
			},
			expectErr: assert.AnError,
		},
		{
			name:   "database error when creating product",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     receptionID,
						PVZID:  pvzID,
						Status: reception_domain.InProgress,
					}, nil)

				p.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			expectErr: assert.AnError,
		},
		{
			name: "invalid params",
			params: application.CreateParams{
				PVZID:    "",
				Type:     product_domain.Electronics,
				UserRole: auth_domain.RoleEmployee,
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {
			},
			expectErr: product_domain.ErrInvalidIDFormat,
		},
		{
			name:   "invalid role params",
			params: invalidRoleParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository, p *product_mocks.MockProductRepository) {

			},
			expectErr: product_domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receptionRepo := reception_mocks.NewMockReceptionRepository(ctrl)
			productRepo := product_mocks.NewMockProductRepository(ctrl)
			tt.mockSetup(receptionRepo, productRepo)

			logger := logger.NewTestLogger()

			service := application.NewProductService(productRepo, receptionRepo, logger)
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
