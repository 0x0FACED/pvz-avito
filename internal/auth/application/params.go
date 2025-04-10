package application

import (
	"errors"
	"fmt"

	"github.com/0x0FACED/pvz-avito/internal/auth/domain"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidRole  = errors.New("invalid role")
)

type RegisterParams struct {
	Email    domain.Email
	Password string
	Role     domain.Role
}

func (p RegisterParams) Validate() error {
	if err := p.Email.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidEmail, err)
	}

	if err := p.Role.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRole, err)
	}

	return nil
}

type LoginParams struct {
	Email    domain.Email
	Password string
}

func (p LoginParams) Validate() error {
	if err := p.Email.Validate(); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	return nil
}
