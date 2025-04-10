package domain

import "context"

type ReceptionRepository interface {
	Create(ctx context.Context, reception *Reception) (*Reception, error)
	FindByID(ctx context.Context, id string) (*Reception, error)
	FindLastOpenByPVZ(ctx context.Context, pvzID string) (*Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (*Reception, error)
	ListByPVZ(ctx context.Context, pvzID string) ([]*Reception, error)
}
