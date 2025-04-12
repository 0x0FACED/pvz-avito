package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"net/url"
	"testing"
	"time"

	product_http "github.com/0x0FACED/pvz-avito/internal/product/delivery/http"
	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	pvz_http "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/http"
	reception_http "github.com/0x0FACED/pvz-avito/internal/reception/delivery/http"
	"github.com/google/uuid"
)

func authUserDummy(t *testing.T, baseURL string, role string) string {
	reqBody := map[string]string{"role": role}
	data, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	addr := fmt.Sprintf("%s/dummyLogin", baseURL)

	resp, err := nethttp.Post(addr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to do dummyLogin request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("unexpected status code from dummyLogin: got %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	token := string(body)

	return token
}

func createPVZ(t *testing.T, baseURL string, token string) string {

	id := uuid.NewString()
	regDate := time.Now()

	reqBody := pvz_http.CreateRequest{
		ID:               &id,
		RegistrationDate: &regDate,
		City:             "Москва",
	}

	data, _ := json.Marshal(reqBody)

	addr, err := url.JoinPath(baseURL, "pvz")
	if err != nil {
		t.Fatalf("failed to join baseURL and pvz: %v", err)
	}

	req, err := nethttp.NewRequest(nethttp.MethodPost, addr, bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to create pvz request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send pvz request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated {
		t.Fatalf("unexpected status code from pvz create: %v", resp.StatusCode)
	}

	var result pvz_http.CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode pvz response: %v", err)
	}

	return *result.ID
}

func createReception(t *testing.T, baseURL string, token string, pvzID string) reception_http.CreateResponse {
	reqBody := reception_http.CreateRequest{
		PVZID: pvzID,
	}

	data, _ := json.Marshal(reqBody)

	addr, err := url.JoinPath(baseURL, "receptions")
	if err != nil {
		t.Fatalf("failed to join baseURL and receptions: %v", err)
	}

	req, err := nethttp.NewRequest(nethttp.MethodPost, addr, bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to create reception request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send create reception request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated {
		t.Fatalf("unexpected status code from reception create: %v", resp.StatusCode)
	}

	var result reception_http.CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode reception response: %v", err)
	}

	return result
}

func createProduct(t *testing.T, baseURL string, token string, pvzID string) product_http.CreateResponse {
	reqBody := product_http.CreateRequest{
		Type:  string(product_domain.Electronics),
		PVZID: pvzID,
	}

	data, _ := json.Marshal(reqBody)

	addr, err := url.JoinPath(baseURL, "products")
	if err != nil {
		t.Fatalf("failed to join baseURL and pvz: %v", err)
	}

	req, err := nethttp.NewRequest(nethttp.MethodPost, addr, bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to create product request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send create product request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated {
		t.Fatalf("unexpected status code from product create: %v", resp.StatusCode)
	}

	var result product_http.CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode product response: %v", err)
	}

	return result
}

func closeReception(t *testing.T, baseURL string, token string, pvzID string) reception_http.CreateResponse {
	addr, err := url.JoinPath(baseURL, "pvz", pvzID, "close_last_reception")
	if err != nil {
		t.Fatalf("failed to join baseURL and pvz: %v", err)
	}

	req, err := nethttp.NewRequest(nethttp.MethodPost, addr, nil)
	if err != nil {
		t.Fatalf("failed to close reception request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send close reception request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("unexpected status code from close reception: %v", resp.StatusCode)
	}

	var result reception_http.CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode close reception response: %v", err)
	}

	return result
}

func createProductAfterCloseReception(t *testing.T, baseURL string, token string, pvzID string) int {
	reqBody := product_http.CreateRequest{
		PVZID: pvzID,
	}

	data, _ := json.Marshal(reqBody)

	addr, err := url.JoinPath(baseURL, "products")
	if err != nil {
		t.Fatalf("failed to join baseURL and pvz: %v", err)
	}

	req, err := nethttp.NewRequest(nethttp.MethodPost, addr, bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to create product request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send create product request: %v", err)
	}

	defer resp.Body.Close()

	return resp.StatusCode
}
