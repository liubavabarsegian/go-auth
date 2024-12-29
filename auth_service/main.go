package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	keycloakClient *gocloak.GoCloak
	realm          = os.Getenv("KEYCLOAK_REALM")
	clientID       = os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret   = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	keycloakURL    = os.Getenv("KEYCLOAK_URL")
)

func main() {
	keycloakClient = gocloak.NewClient(keycloakURL)
	keycloakClient.RestyClient().SetDebug(true)

	log.Println("Attempting to connect to Keycloak at:", keycloakURL)

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	log.Println("Auth service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := keycloakClient.LoginAdmin(r.Context(), "admin", "admin", realm)
	if err != nil {
		log.Println("error in register")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// user := gocloak.User{
	// 	ID:       gocloak.StringP("auth-service"),
	// 	Username: &req.Username,
	// 	Email:    &req.Email,
	// 	Enabled:  gocloak.BoolP(true),
	// 	Credentials: &[]gocloak.CredentialRepresentation{
	// 		{
	// 			Type:      gocloak.StringP("password"),
	// 			Value:     &req.Password,
	// 			Temporary: gocloak.BoolP(false),
	// 		},
	// 	},
	// }

	// // Логируем тело пользователя
	// log.Printf("User payload: %+v", user)
	// log.Printf("Client ID: %s, Client Secret: %s, Keycloak URL: %s, Realm: %s", clientID, clientSecret, keycloakURL, realm)

	// _, err = keycloakClient.CreateUser(r.Context(), token.AccessToken, realm, user)
	// log.Printf("Access Token: %s", token.AccessToken)

	// users, err := keycloakClient.GetUsers(r.Context(), token.AccessToken, realm, gocloak.GetUsersParams{})
	// if err != nil {
	// 	log.Fatalf("Error fetching clients: %v", err)
	// }
	// log.Println(users)

	// if err != nil {
	// 	log.Printf("CreateUser error: %+v", err)
	// 	log.Printf("User payload: %+v", user)
	// 	http.Error(w, "Failed to create user", http.StatusInternalServerError)
	// 	return
	// }

	client := gocloak.Client{
		ClientID:                  gocloak.StringP("auth-service"),
		Secret:                    gocloak.StringP("my-client-secret"), // Ensure this matches the secret generated in Keycloak
		Enabled:                   gocloak.BoolP(true),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
		PublicClient:              gocloak.BoolP(true), // Make sure it's a confidential client
	}

	_, err = keycloakClient.CreateClient(r.Context(), token.AccessToken, realm, client)
	if err != nil {
		log.Println("error in creating client")
		log.Println(err.Error())
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	clients, err := keycloakClient.GetClients(r.Context(), token.AccessToken, realm, gocloak.GetClientsParams{})
	if err != nil {
		log.Fatalf("Error fetching clients: %v", err)
	}
	log.Println(clients)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Use direct access grant for user login
	token, err := keycloakClient.Login(r.Context(), clientID, clientSecret, realm, req.Email, req.Password)
	if err != nil {
		log.Printf("Login error: %+v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
