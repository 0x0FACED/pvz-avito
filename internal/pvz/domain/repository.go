package domain

import (
	"context"
	"time"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz *PVZ) (*PVZ, error)
	FindByID(ctx context.Context, id string) (*PVZ, error)
	ListWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*PVZWithReceptions, error)
}
