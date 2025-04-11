package application

import (
	"errors"
	"time"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	"github.com/google/uuid"
)

type CreateParams struct {
	ID               *string
	RegistrationDate *time.Time
	City             pvz_domain.City
	UserRole         auth_domain.Role
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

	// TODO: refactor
	if p.UserRole != auth_domain.RoleModerator {
		return errors.New("only moderator can create new pvz")
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

type CloseLastReceptionParams struct {
	PVZID    string
	UserRole auth_domain.Role
}

func (p CloseLastReceptionParams) Validate() error {
	if err := uuid.Validate(p.PVZID); err != nil {
		return errors.New("invalid pvz id")
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return errors.New("only employee can close reception")
	}

	return nil
}

type DeleteLastProductParams struct {
	PVZID    string
	UserRole auth_domain.Role
}

func (p DeleteLastProductParams) Validate() error {
	if err := uuid.Validate(p.PVZID); err != nil {
		return errors.New("invalid pvz id")
	}

	if p.UserRole != auth_domain.RoleEmployee {
		return errors.New("only employee can delete products")
	}

	return nil
}
