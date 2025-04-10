package domain

import (
	"time"
)

type ProductType string

const (
	Electronics ProductType = "электроника"
	Clothes     ProductType = "одежда"
	Shoes       ProductType = "обувь"
)

func (p ProductType) String() string {
	return string(p)
}

type Product struct {
	ID          string
	DateTime    time.Time
	Type        ProductType
	ReceptionID string
}
