package application

import "github.com/0x0FACED/pvz-avito/internal/reception/domain"

type ReceptionService struct {
	repo domain.ReceptionRepository
}

func NewReceptionService(repo domain.ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		repo: repo,
	}
}
