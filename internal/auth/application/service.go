package application

import (
	"context"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/google/uuid"
)

type AuthService struct {
	repo auth_domain.UserRepository

	log *logger.ZerologLogger
}

func NewAuthService(repo auth_domain.UserRepository, l *logger.ZerologLogger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  l,
	}
}

func (s *AuthService) Register(ctx context.Context, params RegisterParams) (*auth_domain.User, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("Register")
		return nil, err
	}

	// using min cost, dont use in prod
	// separate cost to .env file
	hash, err := HashPasswordString(params.Password, 4)
	if err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("Error hash password")
		return nil, err
	}

	user := &auth_domain.User{
		ID:       uuid.NewString(),
		Email:    params.Email,
		Password: hash,
		Role:     params.Role,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error().Any("params", params).Any("user", user).Err(err).Msg("Error creating user")
		return nil, err
	}

	s.log.Info().Any("params", params).Any("user", user).Msg("Register successful")

	return created, nil
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (*auth_domain.User, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("Login")
		return nil, err
	}

	user, err := s.repo.FindByEmail(ctx, params.Email.String())
	if err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("Error finding user by email")
		return nil, err
	}

	err = CompareHashAndPassword(user.Password, params.Password)
	if err != nil {
		s.log.Error().Any("params", params).Any("user", user).Err(err).Msg("Password mismatch")
		return nil, err
	}

	s.log.Info().Any("params", params).Any("user", user).Msg("Login successful")
	return user, nil
}
