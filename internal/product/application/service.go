package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
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
		// TODO: handle err correctly
		return nil, err
	}

	lastReception, err := s.receptionRepo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	if lastReception == nil {
		return nil, errors.New("there is no open reception in this pvz")
	}

	product := &product_domain.Product{
		ID:          uuid.NewString(),
		DateTime:    time.Now(),
		Type:        params.Type,
		ReceptionID: lastReception.ID,
	}
	created, err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return created, nil
}
