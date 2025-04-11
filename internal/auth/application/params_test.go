package application_test

import (
	"testing"

	"github.com/0x0FACED/pvz-avito/internal/auth/application"
	"github.com/stretchr/testify/assert"
)

func TestRegisterParams_Validate(t *testing.T) {
	tests := []struct {
		name      string
		params    application.RegisterParams
		expectErr bool
	}{
		{"valid", application.RegisterParams{"test@example.com", "123", "employee"}, false},
		{"empty email", application.RegisterParams{"", "123", "employee"}, true},
		{"invalid role", application.RegisterParams{"test@example.com", "123", "mod"}, true},
		{"invalid email", application.RegisterParams{"test@", "123", "moderator"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoginParams_Validate(t *testing.T) {
	tests := []struct {
		name      string
		params    application.LoginParams
		expectErr bool
	}{
		{"valid", application.LoginParams{"test@example.com", "123"}, false},
		{"empty email", application.LoginParams{"", "123"}, true},
		{"invalid email", application.LoginParams{"test@", "123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
