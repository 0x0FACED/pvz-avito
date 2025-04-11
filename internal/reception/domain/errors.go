package domain

import "errors"

var (
	ErrInternalDatabase   = errors.New("reception: internal database error")
	ErrNoReceptionToClose = errors.New("reception: no reception found")
)
