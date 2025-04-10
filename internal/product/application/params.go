package application

import (
	"errors"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	Type     product_domain.ProductType
	PVZID    string
	UserRole auth_domain.Role
}

func (p CreateParams) Validate() error {
	// TODO: refactor
	if p.Type != product_domain.Shoes && p.Type != product_domain.Electronics && p.Type != product_domain.Clothes {
		return errors.New("unsupported product type")
	}

	if err := uuid.Validate(p.PVZID); err != nil {
		return errors.New("not uuid")
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return errors.New("only employee can add products")
	}

	return nil
}
