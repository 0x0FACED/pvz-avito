package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_PVZ_Reception_Flow(t *testing.T) {
	baseURL := "http://localhost:8080"

	// getting employee and moderator tokens
	moderatorToken := authUserDummy(t, baseURL, "moderator")
	assert.NotEmpty(t, moderatorToken, "moderator token must not be empty")

	employeeToken := authUserDummy(t, baseURL, "employee")
	assert.NotEmpty(t, employeeToken, "employee token must not be empty")

	// creating pvz
	pvzID := createPVZ(t, baseURL, moderatorToken)
	assert.NotEmpty(t, pvzID, "PVZ ID must not be empty")

	// creating reception
	createReceptionResp := createReception(t, baseURL, employeeToken, pvzID)
	assert.NotEmpty(t, createReceptionResp.ID, "Reception ID must not be empty")
	assert.Equal(t, pvzID, createReceptionResp.PVZID, "Reception must be linked to correct PVZ")
	assert.Equal(t, "in_progress", createReceptionResp.Status, "Reception must be initially in_progress")

	productIDs := make(map[string]struct{}, 50)

	// addding 50 products to reception
	for range 50 {
		productResp := createProduct(t, baseURL, employeeToken, pvzID)

		assert.NotEmpty(t, productResp.ID, "Product ID must not be empty")
		assert.Equal(t, createReceptionResp.ID, productResp.ReceptionID, "Product must belong to current reception")
		assert.NotEmpty(t, productResp.Type, "Product Type must not be empty")

		_, exists := productIDs[productResp.ID]
		assert.False(t, exists, "Product ID must be unique")
		productIDs[productResp.ID] = struct{}{}
	}

	// closing reception
	closeReceptionResp := closeReception(t, baseURL, employeeToken, pvzID)

	assert.Equal(t, createReceptionResp.ID, closeReceptionResp.ID, "Reception ID must match")
	assert.Equal(t, pvzID, closeReceptionResp.PVZID, "PVZ ID must match")
	assert.Equal(t, "close", closeReceptionResp.Status, "Reception status must be close")

	// after closing reception trying to add new product to reception
	// must return 400
	respCode := createProductAfterCloseReception(t, baseURL, employeeToken, pvzID)
	assert.Equal(t, http.StatusBadRequest, respCode, "Resp code must be 400 bcz reception is closed")
}
