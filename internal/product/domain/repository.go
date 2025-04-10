package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetLastByReception(ctx context.Context, receptionID string) (*Product, error)
	DeleteLastFromReception(ctx context.Context, receptionID string) error
	ListByReception(ctx context.Context, receptionID string) ([]*Product, error)
}
