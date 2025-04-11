package application

import (
	"fmt"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
)

type RegisterParams struct {
	Email    auth_domain.Email
	Password string
	Role     auth_domain.Role
}

func (p RegisterParams) Validate() error {
	if err := p.Email.Validate(); err != nil {
		return fmt.Errorf("%w: %w", auth_domain.ErrInvalidEmail, err)
	}

	if err := p.Role.Validate(); err != nil {
		return fmt.Errorf("%w: %w", auth_domain.ErrInvalidRole, err)
	}

	return nil
}

type LoginParams struct {
	Email    auth_domain.Email
	Password string
}

func (p LoginParams) Validate() error {
	if err := p.Email.Validate(); err != nil {
		return fmt.Errorf("%w: %w", auth_domain.ErrInvalidEmail, err)
	}

	return nil
}
