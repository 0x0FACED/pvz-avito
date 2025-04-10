package application

import (
	"context"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/google/uuid"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(ctx context.Context, params CreateParams) (*domain.Product, error) {
	product := &domain.Product{
		ID:       uuid.NewString(),
		DateTime: time.Now(),
		Type:     params.Type,
	}
	err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
