package domain

import (
	"fmt"
	"time"

	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
)

type City string

const (
	Moscow City = "Москва"
	SPb    City = "Санкт-Петербург"
	Kazan  City = "Казань"
)

func (c City) String() string {
	return string(c)
}

func (c City) Validate() error {
	if c != Moscow && c != SPb && c != Kazan {
		return fmt.Errorf("unsupported city: %s", c)
	}

	return nil
}

type PVZ struct {
	ID               *string
	RegistrationDate *time.Time
	City             City
}

type PVZWithReceptions struct {
	PVZ        *PVZ
	Receptions []*ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception *reception_domain.Reception
	Products  []*product_domain.Product
}
