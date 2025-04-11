package domain

import "errors"

var (
	// create
	ErrPVZAlreadyExists = errors.New("pvz: pvz already exists")
	ErrInternalDatabase = errors.New("pvz: internal database error")
)

var (
	// create
	ErrUnsupportedCity = errors.New("pvz: unsupported city")
	// create
	ErrInvalidRole = errors.New("pvz: only moderators can create new pvz")
)

var (
	// create
	ErrInvalidIDFormat = errors.New("pvz: invalid id format")
)
