package application

import (
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	ID               *string
	RegistrationDate *time.Time
	City             domain.City
}

func (p CreateParams) Validate() error {
	if p.ID != nil {
		if err := uuid.Validate(*p.ID); err != nil {
			// TODO: handle err
			return err
		}
	}

	if err := p.City.Validate(); err != nil {
		return err
	}

	return nil
}

type ListWithReceptionsParams struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      *int
	Limit     *int
}

func (p *ListWithReceptionsParams) Validate() error {
	if p.Page == nil {
		// default val is 1
		// (incorrect, will panic)
		*p.Page = 1
	} else if *p.Page < 0 {
		return errors.New("page cant be < 0")
	}

	if p.Limit == nil {
		// default val is 1
		// (incorrect, will panic)
		*p.Limit = 10
	} else if *p.Limit < 0 {
		return errors.New("limit cant be < 0")
	}

	return nil
}
