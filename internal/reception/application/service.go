package application

import (
	"context"
	"errors"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/google/uuid"
)

type ReceptionService struct {
	repo reception_domain.ReceptionRepository

	log *logger.ZerologLogger
}

func NewReceptionService(repo reception_domain.ReceptionRepository, l *logger.ZerologLogger) *ReceptionService {
	return &ReceptionService{
		repo: repo,
		log:  l,
	}
}

func (s *ReceptionService) Create(ctx context.Context, params CreateParams) (*reception_domain.Reception, error) {
	if err := params.Validate(); err != nil {
		s.log.Error().Any("params", params).Err(err).Msg("CreateReception")
		return nil, err
	}

	last, err := s.repo.FindLastOpenByPVZ(ctx, params.PVZID)
	if err != nil && !errors.Is(err, reception_domain.ErrNoOpenReception) {
		s.log.Error().Any("params", params).Err(err).Msg("Error checking last open reception")
		return nil, err
	}

	if last != nil {
		s.log.Error().Any("params", params).Err(reception_domain.ErrFoundOpenedReception).Msg("Opened reception already exists")
		return nil, reception_domain.ErrFoundOpenedReception
	}

	reception := reception_domain.Reception{
		ID:       uuid.NewString(),
		DateTime: time.Now(),
		PVZID:    params.PVZID,
		Status:   reception_domain.InProgress,
	}

	created, err := s.repo.Create(ctx, &reception)
	if err != nil {
		s.log.Error().Any("params", params).Any("reception", reception).Err(err).Msg("Error creating reception")
		return nil, err
	}

	s.log.Info().Any("params", params).Any("reception", created).Msg("CreateReception successful")
	return created, nil
}
