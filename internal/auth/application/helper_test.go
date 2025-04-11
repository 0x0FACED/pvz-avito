package application_test

import (
	"testing"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_HashPasswordString(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		cost      int
		expectErr bool
	}{
		{"valid password", "secure123", bcrypt.DefaultCost, false},
		{"empty password", "", bcrypt.DefaultCost, false},
		{"too high cost", "secure123", 1000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := application.HashPasswordString(tt.password, tt.cost)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func Test_CompareHashAndPassword(t *testing.T) {
	password := "test123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		hash      string
		password  string
		expectErr bool
	}{
		{"correct password", string(hash), password, false},
		{"wrong password", string(hash), "wrong", true},
		{"invalid hash", "notahash", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := application.CompareHashAndPassword(tt.hash, tt.password)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
