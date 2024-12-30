package main

import (
	"authService/router"
	"log"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

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

	mux := router.SetUpRouter(keycloakClient)
	log.Println("Attempting to connect to Keycloak at:", keycloakURL)

	log.Println("Auth service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
