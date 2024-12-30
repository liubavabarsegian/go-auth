package login

import (
	"authService/config"
	"encoding/json"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

func Login(w http.ResponseWriter, r *http.Request, keycloakClient *gocloak.GoCloak) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(loginErrorResponse{Status: http.StatusMethodNotAllowed, Error: "Method not allowed"})
		return
	}

	var request loginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(loginErrorResponse{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	token, err := keycloakClient.Login(
		r.Context(),
		config.ClientID,
		config.ClientSecret,
		config.Realm,
		request.Username,
		request.Password,
	)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(loginErrorResponse{Status: http.StatusUnauthorized, Error: err.Error()})
		return
	}

	resp := loginResponse{
		Status:       http.StatusOK,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
