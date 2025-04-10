package application

import (
	"errors"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	PVZID    string
	UserRole auth_domain.Role
}

func (p CreateParams) Validate() error {
	if err := uuid.Validate(p.PVZID); err != nil {
		// TODO: handle err
		return err
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return errors.New("only employee can create reception")
	}

	return nil
}
