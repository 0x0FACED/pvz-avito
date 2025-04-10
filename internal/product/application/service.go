package application

import (
	"context"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/google/uuid"
)

type ProductService struct {
	repo domain.ProductRepository

	log *logger.ZerologLogger
}

func NewProductService(repo domain.ProductRepository, l *logger.ZerologLogger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  l,
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
