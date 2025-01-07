package logout

import (
	"authService/config"
	"encoding/json"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

func Logout(w http.ResponseWriter, r *http.Request, keycloakClient *gocloak.GoCloak) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(logoutErrorResponse{Status: http.StatusMethodNotAllowed, Error: "Method not allowed"})
		return
	}

	var request logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(logoutErrorResponse{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	// Check realm value
	if config.Realm == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(logoutErrorResponse{Status: http.StatusInternalServerError, Error: "Realm not set"})
		return
	}

	err := keycloakClient.Logout(r.Context(), config.ClientID, config.ClientSecret, config.Realm, request.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(logoutErrorResponse{Status: http.StatusInternalServerError, Error: "Failed to logout user: " + err.Error()})
		return
	}

	resp := logoutResponse{
		Status:  http.StatusOK,
		Message: "User logged out successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
