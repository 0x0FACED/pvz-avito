package http

import (
	"time"
)

type CreateResponse struct {
	ID               *string    `json:"id,omitempty"`
	RegistrationDate *time.Time `json:"registrationDate,omitempty"`
	City             string     `json:"city"`
}

type CloseResponse struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type ListResponse struct {
	PVZ        pvz         `json:"pvz"`
	Receptions []reception `json:"receptions"`
}

type product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

type reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
	Products []product `json:"products"`
}

type pvz struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}
