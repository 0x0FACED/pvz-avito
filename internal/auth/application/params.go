package application

import (
	"fmt"

	"github.com/0x0FACED/pvz-avito/internal/auth/domain"
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
		return fmt.Errorf("%w: %w", ErrInvalidEmail, err)
	}

	return nil
}
