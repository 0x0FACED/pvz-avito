package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("auth: user already exists")
	ErrInternalDatabase  = errors.New("auth: internal database error")
	ErrUserNotFound      = errors.New("auth: user not found")
)

var (
	ErrHashPassword    = errors.New("auth: invalid password format")
	ErrInvalidPassword = errors.New("auth: invalid password")
)

var (
	ErrInvalidEmail = errors.New("auth: invalid email")
	ErrInvalidRole  = errors.New("auth: invalid role")
)
