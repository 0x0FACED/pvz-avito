package domain

import "errors"

var (
	ErrInternalDatabase     = errors.New("reception: internal database error")
	ErrReceptionNotFound    = errors.New("reception: reception not found")
	ErrNoOpenReception      = errors.New("reception: no open reception found")
	ErrPVZNotFound          = errors.New("reception: pvz not found")
	ErrFoundOpenedReception = errors.New("reception: there is opened reception, cant create new one")
)

var (
	ErrAccessDenied = errors.New("reception: only employees can create new reception")
)

var (
	ErrInvalidIDFormat = errors.New("reception: invalid id format")
)
