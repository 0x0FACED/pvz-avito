package application

import (
	"errors"

	"github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	Type  domain.ProductType
	PVZID string
}

func (p CreateParams) Validate() error {
	// TODO: refactor
	if p.Type != domain.Shoes && p.Type != domain.Electronics && p.Type != domain.Clothes {
		return errors.New("unsupported product type")
	}

	if err := uuid.Validate(p.PVZID); err != nil {
		return errors.New("not uuid")
	}

	return nil
}
