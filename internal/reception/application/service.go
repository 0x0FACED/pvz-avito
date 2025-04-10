package application

import (
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/reception/domain"
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
