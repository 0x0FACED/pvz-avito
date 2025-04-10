package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
)

type ReceptionService struct {
	repo domain.ReceptionRepository

	log *logger.ZerologLogger
}

func NewReceptionService(repo domain.ReceptionRepository, l *logger.ZerologLogger) *ReceptionService {
	return &ReceptionService{
		repo: repo,
		log:  l,
	}
}

func (s *ReceptionService) Create(ctx context.Context, params CreateParams) (*domain.Reception, error) {
	last, err := s.repo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	// TODO: refactor
	if last != nil {
		// TODO: handle err
		return nil, errors.New("there is in progress reception, cant create new one")
	}

	// all is good, can create new reception
	reception := domain.Reception{
		ID:       uuid.NewString(),
		DateTime: time.Now(),
		PVZID:    params.PVZID,
		Status:   domain.InProgress,
	}

	created, err := s.repo.Create(ctx, &reception)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return created, nil
}
