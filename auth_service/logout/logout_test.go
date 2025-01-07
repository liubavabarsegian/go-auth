package logout

import (
	"authService/config"
	"bytes"
	"context"
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

func TestLogoutSuccess(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	token, err := client.Login(context.Background(), config.ClientID, config.ClientSecret, config.Realm, "admin", "admin")
	if err != nil {
		t.Fatalf("Failed to log in: %v", err)
	}

	requestBody, _ := json.Marshal(logoutRequest{RefreshToken: token.RefreshToken})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Logout(rec, req, client)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp logoutResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "User logged out successfully", resp.Message)
}

func TestLogoutInvalidRefreshToken(t *testing.T) {
	setTestEnvVars()
	config.Realm = os.Getenv("KEYCLOAK_REALM")
	config.ClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	config.ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	config.KeycloakURL = os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(config.KeycloakURL)

	// Use an invalid refresh token for testing
	requestBody, _ := json.Marshal(logoutRequest{RefreshToken: "invalid_refresh_token"})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Logout(rec, req, client)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp logoutErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Error, "Failed to logout user")
}

func TestLogoutInvalidMethod(t *testing.T) {
	setTestEnvVars()
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Logout(rec, req, client)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	var resp logoutErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Method not allowed", resp.Error)
}

func TestLogoutRealmNotSet(t *testing.T) {
	setTestEnvVars()
	config.Realm = "" // Simulate unset realm

	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))

	// Use a valid refresh token for testing
	requestBody, _ := json.Marshal(logoutRequest{RefreshToken: "valid_refresh_token"})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Logout(rec, req, client)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp logoutErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Error, "Realm not set")
}
