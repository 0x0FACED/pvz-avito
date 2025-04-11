package http

import "time"

type createRequest struct {
	ID               *string    `json:"id,omitempty"`
	RegistrationDate *time.Time `json:"registrationDate,omitempty"`
	City             string     `json:"city"`
}
