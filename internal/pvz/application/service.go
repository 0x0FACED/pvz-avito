package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/pkg/metrics"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
)

type PVZService struct {
	pvzRepo       pvz_domain.PVZRepository
	receptionRepo reception_domain.ReceptionRepository
	productRepo   product_domain.ProductRepository

	log *logger.ZerologLogger
}

func NewPVZService(
	pvzRepo pvz_domain.PVZRepository,
	receptionRepo reception_domain.ReceptionRepository,
	productRepo product_domain.ProductRepository,
	l *logger.ZerologLogger,
) *PVZService {
	return &PVZService{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		log:           l,
	}
}

func (s *PVZService) Create(ctx context.Context, params CreateParams) (*pvz_domain.PVZ, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("CreatePVZ")
		return nil, err
	}

	pvz := pvz_domain.PVZ{
		ID:               params.ID,
		RegistrationDate: params.RegistrationDate,
		City:             params.City,
	}

	if pvz.ID == nil {
		id := uuid.NewString()
		pvz.ID = &id
	}
	if pvz.RegistrationDate == nil {
		now := time.Now()
		pvz.RegistrationDate = &now
	}

	created, err := s.pvzRepo.Create(ctx, &pvz)
	if err != nil {
		s.log.Error().Any("params", params).Any("pvz", pvz).Err(err).Msg("Error creating PVZ")
		return nil, err
	}

	metrics.PvzCreatedTotal.Inc()

	s.log.Info().Any("params", params).Any("pvz", created).Msg("CreatePVZ successful")

	return created, nil
}

func (s *PVZService) CloseLastReception(ctx context.Context, params CloseLastReceptionParams) (*reception_domain.Reception, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("CloseLastReception")
		return nil, err
	}

	reception, err := s.receptionRepo.CloseLastReception(ctx, params.PVZID)
	if err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("Error closing last reception")
		return nil, err
	}

	s.log.Info().Any("params", params).Any("reception", reception).Msg("CloseLastReception successful")

	return reception, nil
}

func (s *PVZService) DeleteLastProduct(ctx context.Context, params DeleteLastProductParams) error {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("DeleteLastProduct")
		return err
	}

	reception, err := s.receptionRepo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil && !errors.Is(err, reception_domain.ErrNoOpenReception) {
		s.log.Error().Any("params", params).Err(err).Msg("Error finding last open reception")
		return err
	}

	if reception == nil {
		s.log.Error().Any("params", params).Err(reception_domain.ErrNoOpenReception).Msg("No open reception")
		return reception_domain.ErrNoOpenReception
	}

	err = s.productRepo.DeleteLastFromReception(ctx, reception.ID)
	if err != nil {
		s.log.Error().Any("params", params).Any("reception", reception).Err(err).Msg("Error deleting last product")
		return err
	}

	s.log.Info().Any("params", params).Any("reception", reception).Msg("DeleteLastProduct successful")
	return nil
}

func (s *PVZService) ListAllPVZs(ctx context.Context) ([]*pvz_domain.PVZ, error) {
	pvzs, err := s.pvzRepo.ListAllPVZs(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("ListAllPVZs")
		return pvzs, err
	}

	return pvzs, nil
}

func (s *PVZService) ListWithReceptions(ctx context.Context, params ListWithReceptionsParams) ([]*pvz_domain.PVZWithReceptions, error) {
	pvzWithReceptions, err := s.pvzRepo.ListWithReceptions(ctx, params.StartDate, params.EndDate, *params.Page, *params.Limit)
	if err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("ListWithReceptions")
		return nil, err
	}

	s.log.Info().Any("params", params).Any("resultCount", len(pvzWithReceptions)).Msg("ListWithReceptions successful")

	return pvzWithReceptions, nil
}
