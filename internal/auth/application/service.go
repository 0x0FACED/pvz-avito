package application

import (
	"context"
	"errors"

	"github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo domain.UserRepository

	log *logger.ZerologLogger
}

func NewAuthService(repo domain.UserRepository, l *logger.ZerologLogger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  l,
	}
}

func (s *AuthService) Register(ctx context.Context, params RegisterParams) (*domain.User, error) {
	// using min cost, dont use in prod
	// separate cost to .env file
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	user := &domain.User{
		ID:       uuid.NewString(),
		Email:    params.Email,
		Password: string(hash),
		Role:     params.Role,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (*domain.User, error) {
	if err := params.Validate(); err != nil {
		// TODO: handle err
		return nil, err
	}

	user, err := s.repo.FindByEmail(ctx, params.Email.String())
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		// TODO: handle err
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
