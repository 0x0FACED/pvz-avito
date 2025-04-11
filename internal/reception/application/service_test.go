package application_test

import (
	"context"
	"testing"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/reception/application"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	reception_mocks "github.com/0x0FACED/pvz-avito/internal/reception/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReceptionCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := uuid.NewString()

	validParams := application.CreateParams{
		PVZID:    pvzID,
		UserRole: auth_domain.RoleEmployee,
	}

	invalidRoleParams := application.CreateParams{
		PVZID:    pvzID,
		UserRole: auth_domain.RoleModerator,
	}

	tests := []struct {
		name      string
		params    application.CreateParams
		mockSetup func(*reception_mocks.MockReceptionRepository)
		expectErr error
	}{
		{
			name:   "successful creation",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, reception_domain.ErrNoOpenReception)

				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec *reception_domain.Reception) (*reception_domain.Reception, error) {
						assert.Equal(t, pvzID, rec.PVZID)
						assert.Equal(t, reception_domain.InProgress, rec.Status)
						return rec, nil
					})
			},
			expectErr: nil,
		},
		{
			name:   "open reception already exists",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(&reception_domain.Reception{
						ID:     uuid.NewString(),
						PVZID:  pvzID,
						Status: reception_domain.InProgress,
					}, nil)
			},
			expectErr: reception_domain.ErrFoundOpenedReception,
		},
		{
			name:   "database error when checking open reception",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, assert.AnError)
			},
			expectErr: assert.AnError,
		},
		{
			name:   "database error when creating reception",
			params: validParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
				r.EXPECT().
					FindLastOpenByPVZ(gomock.Any(), pvzID).
					Return(nil, reception_domain.ErrNoOpenReception)

				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			expectErr: assert.AnError,
		},
		{
			name: "invalid params",
			params: application.CreateParams{
				PVZID: "",
			},
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
			},
			expectErr: reception_domain.ErrInvalidIDFormat,
		},
		{
			name:   "invalid role params",
			params: invalidRoleParams,
			mockSetup: func(r *reception_mocks.MockReceptionRepository) {
			},
			expectErr: reception_domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := reception_mocks.NewMockReceptionRepository(ctrl)
			tt.mockSetup(repo)

			logger := logger.NewTestLogger()

			service := application.NewReceptionService(repo, logger)
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
