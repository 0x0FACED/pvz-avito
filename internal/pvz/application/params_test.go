package application_test

import (
	"testing"
	"time"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_PVZCreateParams_Validate(t *testing.T) {
	validID := uuid.New().String()
	invalidID := "notanuuid"
	now := time.Now()

	tests := []struct {
		name      string
		params    application.CreateParams
		expectErr bool
	}{
		{"valid", application.CreateParams{ID: &validID, RegistrationDate: &now, City: pvz_domain.City("Москва"), UserRole: auth_domain.RoleModerator}, false},
		{"valid 2", application.CreateParams{ID: &validID, RegistrationDate: &now, City: pvz_domain.City("Казань"), UserRole: auth_domain.RoleModerator}, false},
		{"invalid UUID", application.CreateParams{ID: &invalidID, City: pvz_domain.City("Москва"), UserRole: auth_domain.RoleModerator}, true},
		{"access denied", application.CreateParams{ID: &validID, City: pvz_domain.City("Москва"), UserRole: auth_domain.RoleEmployee}, true},
		{"invalid city", application.CreateParams{ID: &validID, City: pvz_domain.City("Нью-Йорк"), UserRole: auth_domain.RoleModerator}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}

func Test_ListWithReceptionsParams_Validate(t *testing.T) {
	t.Run("nil page and limit", func(t *testing.T) {
		page, limit := (*int)(nil), (*int)(nil)
		p := &application.ListWithReceptionsParams{Page: page, Limit: limit}
		assert.Panics(t, func() { _ = p.Validate() })
	})

	t.Run("negative page", func(t *testing.T) {
		page := -1
		limit := 10
		p := &application.ListWithReceptionsParams{Page: &page, Limit: &limit}
		assert.Error(t, p.Validate())
	})

	t.Run("valid", func(t *testing.T) {
		page := 1
		limit := 10
		p := &application.ListWithReceptionsParams{Page: &page, Limit: &limit}
		assert.NoError(t, p.Validate())
	})
}

func Test_CloseLastReceptionParams_Validate(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name      string
		params    application.CloseLastReceptionParams
		expectErr bool
	}{
		{"valid", application.CloseLastReceptionParams{PVZID: validID, UserRole: auth_domain.RoleEmployee}, false},
		{"invalid UUID", application.CloseLastReceptionParams{PVZID: "invalid", UserRole: auth_domain.RoleEmployee}, true},
		{"access denied", application.CloseLastReceptionParams{PVZID: validID, UserRole: auth_domain.RoleModerator}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}

func Test_DeleteLastProductParams_Validate(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name      string
		params    application.DeleteLastProductParams
		expectErr bool
	}{
		{"valid", application.DeleteLastProductParams{PVZID: validID, UserRole: auth_domain.RoleEmployee}, false},
		{"invalid UUID", application.DeleteLastProductParams{PVZID: "notanuuid", UserRole: auth_domain.RoleEmployee}, true},
		{"access denied", application.DeleteLastProductParams{PVZID: validID, UserRole: auth_domain.RoleModerator}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}
