package domain

import "errors"

var (
	ErrInternalDatabase  = errors.New("reception: internal database error")
	ErrReceptionNotFound = errors.New("reception: reception not found")
	ErrNoOpenReception   = errors.New("reception: no open reception found")
	ErrPVZNotFound       = errors.New("reception: pvz not found")
)
