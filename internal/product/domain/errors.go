package domain

import "errors"

var (
	ErrProductNotFound  = errors.New("product: product not found")
	ErrInternalDatabase = errors.New("product: internal database error")

	ErrReceptionNotFound  = errors.New("product: reception not found")
	ErrNoProductsToDelete = errors.New("product: no products to delete")

	ErrInvalidProductType = errors.New("product: invalid product type")
	ErrInvalidIDFormat    = errors.New("product: invalid id format")
	ErrAccessDenied       = errors.New("product: only employee can add new products")
)
