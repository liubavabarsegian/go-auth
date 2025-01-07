package login

import (
	"authService/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
)

func setTestEnvVars() {
	_ = os.Setenv("KEYCLOAK_REALM", "master")
	_ = os.Setenv("KEYCLOAK_CLIENT_ID", "admin-cli")
	_ = os.Setenv("KEYCLOAK_CLIENT_SECRET", "admin")
	_ = os.Setenv("KEYCLOAK_URL", "http://localhost:8080")
}

func TestLoginIntegration(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	// Use valid test user credentials
	requestBody, _ := json.Marshal(loginRequest{Username: "Test", Password: "TestPassword"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Login(rec, req, client)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "401 Unauthorized")
}

func TestLoginInvalidCredentials(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	// Use invalid credentials
	requestBody, _ := json.Marshal(loginRequest{Username: "invalid_user", Password: "wrong_password"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Login(rec, req, client)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "401 Unauthorized")
	assert.Contains(t, resp["error"], "invalid_grant: Invalid user credentials")
}

func TestLoginMissingFields(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	// Missing password
	requestBody, _ := json.Marshal(loginRequest{Username: "testuser"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Login(rec, req, client)

	assert.Equal(t, http.StatusUnauthorized, rec.Code) // Expect 401 for missing fields.

	var resp map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "401 Unauthorized")
	assert.Contains(t, resp["error"], "invalid_grant: Invalid user credentials")
}

func TestLoginInvalidMethod(t *testing.T) {
	setTestEnvVars()
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Login(rec, req, client)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Method not allowed", resp["error"])
}
