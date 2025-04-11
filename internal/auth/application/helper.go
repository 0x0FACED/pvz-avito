package application

import (
	"fmt"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

func HashPasswordString(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("%w: %w", auth_domain.ErrHashPassword, err)
	}

	return string(hash), nil
}

func CompareHashAndPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("%w: %w", auth_domain.ErrInvalidPassword, err)
	}

	return nil
}
