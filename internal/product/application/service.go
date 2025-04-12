package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/pkg/metrics"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo   product_domain.ProductRepository
	receptionRepo reception_domain.ReceptionRepository

	log *logger.ZerologLogger
}

func NewProductService(productRepo product_domain.ProductRepository, receptionRepo reception_domain.ReceptionRepository, l *logger.ZerologLogger) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
		log:           l,
	}
}

func (s *ProductService) Create(ctx context.Context, params CreateParams) (*product_domain.Product, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("CreateProduct")
		return nil, err
	}

	lastReception, err := s.receptionRepo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil && !errors.Is(err, reception_domain.ErrNoOpenReception) {
		s.log.Error().Any("params", params).Err(err).Msg("Error finding last open reception")
		return nil, err
	}

	if lastReception == nil {
		s.log.Error().Any("params", params).Err(err).Msg("No open reception found")
		return nil, reception_domain.ErrNoOpenReception
	}

	product := &product_domain.Product{
		ID:          uuid.NewString(),
		DateTime:    time.Now(),
		Type:        params.Type,
		ReceptionID: lastReception.ID,
	}

	created, err := s.productRepo.Create(ctx, product)
	if err != nil {
		s.log.Error().Any("params", params).Any("product", product).Err(err).Msg("Error creating product")
		return nil, err
	}

	metrics.ProductsAddedTotal.Inc()

	s.log.Info().Any("params", params).Any("product", created).Msg("CreateProduct successful")
	return created, nil
}
