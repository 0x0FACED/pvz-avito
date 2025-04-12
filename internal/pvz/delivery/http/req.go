package http

import "time"

type CreateRequest struct {
	ID               *string    `json:"id,omitempty"`
	RegistrationDate *time.Time `json:"registrationDate,omitempty"`
	City             string     `json:"city"`
}
