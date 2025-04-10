package application

import (
	"context"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
)

type PVZService struct {
	pvzRepo       pvz_domain.PVZRepository
	receptionRepo reception_domain.ReceptionRepository

	log *logger.ZerologLogger
}

func NewPVZService(pvzRepo pvz_domain.PVZRepository, receptionRepo reception_domain.ReceptionRepository, l *logger.ZerologLogger) *PVZService {
	return &PVZService{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		log:           l,
	}
}

func (s *PVZService) Create(ctx context.Context, params CreateParams) (*pvz_domain.PVZ, error) {
	if err := params.Validate(); err != nil {
		// TODO: handle err correctly
		return nil, err
	}

	pvz := pvz_domain.PVZ{}
	if params.ID != nil {
		pvz.ID = *params.ID
	}

	if params.RegistrationDate != nil {
		pvz.RegistrationDate = *params.RegistrationDate
	}

	pvz.City = params.City

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

func (s *PVZService) ListWithReceptions(ctx context.Context, params ListWithReceptionsParams) ([]*pvz_domain.PVZWithReceptions, error) {
	pvzWithReceptions, err := s.pvzRepo.ListWithReceptions(ctx, params.StartDate, params.EndDate, *params.Page, *params.Limit)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return pvzWithReceptions, nil
}
