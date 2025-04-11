package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
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
		// TODO: handle err correctly
		return nil, err
	}

	pvz := pvz_domain.PVZ{
		ID:               params.ID,
		RegistrationDate: params.RegistrationDate,
		City:             params.City,
	}

	if pvz.ID == nil {
		uuidStr := uuid.NewString()
		pvz.ID = &uuidStr
	}
	if pvz.RegistrationDate == nil {
		now := time.Now()
		pvz.RegistrationDate = &now
	}

	created, err := s.pvzRepo.Create(ctx, &pvz)
	if err != nil {
		// TOOD: handle err
		return nil, err
	}

	return created, nil
}

func (s *PVZService) FindByID(ctx context.Context, id string) (*pvz_domain.PVZ, error) {
	pvz, err := s.pvzRepo.FindByID(ctx, id)
	if err != nil {
		// TOOD: handle err
		return nil, err
	}

	return pvz, nil
}

func (s *PVZService) CloseLastReception(ctx context.Context, params CloseLastReceptionParams) (*reception_domain.Reception, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	reception, err := s.receptionRepo.CloseLastReception(ctx, params.PVZID)
	if err != nil {
		return nil, err
	}

	return reception, nil
}

func (s *PVZService) DeleteLastProduct(ctx context.Context, params DeleteLastProductParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	reception, err := s.receptionRepo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil {
		// TODO: handle err
		return err
	}

	if reception == nil {
		return errors.New("reception is closed")
	}

	err = s.productRepo.DeleteLastFromReception(ctx, reception.ID)
	if err != nil {
		// TODO: handle err
		return err
	}

	return nil
}

func (s *PVZService) ListWithReceptions(ctx context.Context, params ListWithReceptionsParams) ([]*pvz_domain.PVZWithReceptions, error) {
	pvzWithReceptions, err := s.pvzRepo.ListWithReceptions(ctx, params.StartDate, params.EndDate, *params.Page, *params.Limit)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return pvzWithReceptions, nil
}
