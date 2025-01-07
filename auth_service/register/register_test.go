package register

import (
	"authService/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
)

// Helper function to set environment variables for tests
func setTestEnvVars() {
	_ = os.Setenv("KEYCLOAK_REALM", "auth-service-realm")
	_ = os.Setenv("KEYCLOAK_CLIENT_ID", "auth-service")
	_ = os.Setenv("KEYCLOAK_CLIENT_SECRET", "my-client-secret")
	_ = os.Setenv("KEYCLOAK_URL", "http://localhost:8080")
}

func TestRegisterSuccess(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	uniqueUsername := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	requestBody, _ := json.Marshal(RegisterRequest{
		Username: uniqueUsername,
		Email:    uniqueUsername + "@example.com",
		Password: "TestPassword123",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Register(rec, req, client)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp registerResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", resp.Message)
}

func TestRegisterInvalidMethod(t *testing.T) {
	setTestEnvVars()
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))

	// Attempt to register with GET method
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Register(rec, req, client)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	var resp registerErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Method not allowed", resp.Error)
}

func TestRegisterDuplicateUser(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	// User already exists
	requestBody, _ := json.Marshal(RegisterRequest{
		Username: "existinguser",
		Email:    "existinguser@example.com",
		Password: "securepassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Register(rec, req, client)
	Register(rec, req, client)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp registerErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Error, "User exists with same username")
}
