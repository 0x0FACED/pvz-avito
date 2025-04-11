package application

import (
	"fmt"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	PVZID    string
	UserRole auth_domain.Role
}

func (p CreateParams) Validate() error {
	if err := uuid.Validate(p.PVZID); err != nil {
		return fmt.Errorf("%w: %w", reception_domain.ErrInvalidIDFormat, err)
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return reception_domain.ErrAccessDenied
	}

	return nil
}
