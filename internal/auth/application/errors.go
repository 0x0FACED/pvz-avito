package application

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrHashPassword      = errors.New("invalid password format")
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidRole  = errors.New("invalid role")
)
