package register

import (
	"authService/config"
	"encoding/json"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

func Register(w http.ResponseWriter, r *http.Request, keycloakClient *gocloak.GoCloak) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(registerErrorResponse{Status: http.StatusMethodNotAllowed, Error: "Method not allowed"})
		return
	}

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(registerErrorResponse{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	token, err := keycloakClient.LoginAdmin(r.Context(), "admin", "admin", "master")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(registerErrorResponse{Status: http.StatusInternalServerError, Error: "Method not allowed"})
		return
	}

	user := gocloak.User{
		Username: &req.Username,
		Email:    &req.Email,
		Enabled:  gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     &req.Password,
				Temporary: gocloak.BoolP(false),
			},
		},
	}

	_, err = keycloakClient.CreateUser(r.Context(), token.AccessToken, config.Realm, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(registerErrorResponse{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	resp := registerResponse{
		Status:  http.StatusCreated,
		Message: "User registered successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
