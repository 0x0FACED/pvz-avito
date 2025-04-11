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
		return nil, err
	}

	// using min cost, dont use in prod
	// separate cost to .env file
	hash, err := hashPasswordString(params.Password, 4)
	if err != nil {
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
		return nil, err
	}

	return created, nil
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (*auth_domain.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByEmail(ctx, params.Email.String())
	if err != nil {
		return nil, err
	}

	err = compareHashAndPassword(user.Password, params.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
