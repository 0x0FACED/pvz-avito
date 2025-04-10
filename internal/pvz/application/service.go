package application

import (
	"context"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/pvz/domain"
)

type PVZService struct {
	repo domain.PVZRepository

	log *logger.ZerologLogger
}

func NewPVZService(repo domain.PVZRepository, l *logger.ZerologLogger) *PVZService {
	return &PVZService{
		repo: repo,
		log:  l,
	}
}

func (s *PVZService) Create(ctx context.Context, params CreateParams) (*domain.PVZ, error) {
	pvz := domain.PVZ{}
	if params.ID != nil {
		pvz.ID = *params.ID
	}

	if params.RegistrationDate != nil {
		pvz.RegistrationDate = *params.RegistrationDate
	}

	pvz.City = params.City

	created, err := s.repo.Create(ctx, &pvz)
	if err != nil {
		// TOOD: handle err
		return nil, err
	}

	return created, nil
}

func (s *PVZService) FindByID(ctx context.Context, id string) (*domain.PVZ, error) {
	pvz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// TOOD: handle err
		return nil, err
	}

	return pvz, nil
}

func (s *PVZService) ListWithReceptions(ctx context.Context, params ListWithReceptionsParams) ([]*domain.PVZWithReceptions, error) {
	pvzWithReceptions, err := s.repo.ListWithReceptions(ctx, params.StartDate, params.EndDate, *params.Page, *params.Limit)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return pvzWithReceptions, nil
}
