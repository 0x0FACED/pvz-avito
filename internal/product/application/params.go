package application

import (
	"fmt"

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
	if p.Type != product_domain.Shoes && p.Type != product_domain.Electronics && p.Type != product_domain.Clothes {
		return fmt.Errorf("%w: %s", product_domain.ErrInvalidProductType, p.Type)
	}

	if err := uuid.Validate(p.PVZID); err != nil {
		return fmt.Errorf("%w: %w", product_domain.ErrInvalidIDFormat, err)
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return product_domain.ErrAccessDenied
	}

	return nil
}
