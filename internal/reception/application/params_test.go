package application_test

import (
	"testing"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/reception/application"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ReceptionCreateParams_Validate(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name    string
		params  application.CreateParams
		wantErr bool
	}{
		{"valid", application.CreateParams{PVZID: validID, UserRole: auth_domain.RoleEmployee}, false},
		{"invalid UUID", application.CreateParams{PVZID: "notanuuid", UserRole: auth_domain.RoleEmployee}, true},
		{"access denied", application.CreateParams{PVZID: validID, UserRole: auth_domain.RoleModerator}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
