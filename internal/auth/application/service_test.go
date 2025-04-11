package application_test

import (
	"context"
	"testing"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/auth/mocks"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validParams := application.RegisterParams{
		Email:    "test@example.com",
		Password: "securePassword",
		Role:     auth_domain.RoleModerator,
	}

	existsParams := application.RegisterParams{
		Email:    "busymail@example.com",
		Password: "securePassword",
		Role:     auth_domain.RoleModerator,
	}

	errMailParams := application.RegisterParams{
		Email:    "errmail@example.com",
		Password: "securePassword",
		Role:     auth_domain.RoleEmployee,
	}

	tests := []struct {
		name      string
		params    application.RegisterParams
		mockSetup func(m *mocks.MockUserRepository)
		expectErr error
	}{
		{
			name:   "successful register",
			params: validParams,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, u *auth_domain.User) (*auth_domain.User, error) {
						return u, nil
					})
			},
			expectErr: nil,
		},
		{
			name:   "user already exists",
			params: existsParams,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrUserAlreadyExists)
			},
			expectErr: auth_domain.ErrUserAlreadyExists,
		},
		{
			name:   "failed check by mail (db err)",
			params: errMailParams,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, auth_domain.ErrInternalDatabase)
			},
			expectErr: auth_domain.ErrInternalDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			tt.mockSetup(mockRepo)

			logger := logger.NewTestLogger()

			svc := application.NewAuthService(mockRepo, logger)
			_, err := svc.Register(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validHash, err := application.HashPasswordString("validPass123", 4)
	require.NoError(t, err)

	tests := []struct {
		name      string
		params    application.LoginParams
		mockSetup func(*mocks.MockUserRepository)
		expectErr error
	}{
		{
			name: "successful login",
			params: application.LoginParams{
				Email:    "user@example.com",
				Password: "validPass123",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByEmail(gomock.Any(), "user@example.com").
					Return(&auth_domain.User{
						Email:    "user@example.com",
						Password: validHash,
					}, nil)
			},
			expectErr: nil,
		},
		{
			name: "user not found",
			params: application.LoginParams{
				Email:    "missing@example.com",
				Password: "anyPass",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByEmail(gomock.Any(), "missing@example.com").
					Return(nil, auth_domain.ErrUserNotFound)
			},
			expectErr: auth_domain.ErrUserNotFound,
		},
		{
			name: "invalid password",
			params: application.LoginParams{
				Email:    "user@example.com",
				Password: "wrongPass",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByEmail(gomock.Any(), "user@example.com").
					Return(&auth_domain.User{
						Email:    "user@example.com",
						Password: validHash,
					}, nil)
			},
			expectErr: auth_domain.ErrInvalidPassword,
		},
		{
			name: "database error",
			params: application.LoginParams{
				Email:    "dberror@example.com",
				Password: "anyPass",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByEmail(gomock.Any(), "dberror@example.com").
					Return(nil, auth_domain.ErrInternalDatabase)
			},
			expectErr: auth_domain.ErrInternalDatabase,
		},
		{
			name: "invalid email format",
			params: application.LoginParams{
				Email:    "invalid_email",
				Password: "anyPass",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
			},
			expectErr: auth_domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			tt.mockSetup(mockRepo)

			log := logger.NewTestLogger()
			service := application.NewAuthService(mockRepo, log)

			_, err := service.Login(context.Background(), tt.params)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
