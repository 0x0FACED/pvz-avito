package application_test

import (
	"testing"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/product/application"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ProductCreateParams_Validate(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name      string
		params    application.CreateParams
		expectErr bool
	}{
		{"valid", application.CreateParams{Type: product_domain.Shoes, PVZID: validID, UserRole: auth_domain.RoleEmployee}, false},
		{"invalid type", application.CreateParams{Type: "Food", PVZID: validID, UserRole: auth_domain.RoleEmployee}, true},
		{"invalid UUID", application.CreateParams{Type: product_domain.Shoes, PVZID: "bad-uuid", UserRole: auth_domain.RoleEmployee}, true},
		{"access denied", application.CreateParams{Type: product_domain.Shoes, PVZID: validID, UserRole: auth_domain.RoleModerator}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}
