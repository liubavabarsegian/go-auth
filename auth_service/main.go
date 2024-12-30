package main

import (
	"authService/config"
	"authService/router"
	"log"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

var keycloakClient *gocloak.GoCloak

func main() {
	keycloakClient = gocloak.NewClient(config.KeycloakURL)
	keycloakClient.RestyClient().SetDebug(true)

	mux := router.SetUpRouter(keycloakClient)
	log.Println("Attempting to connect to Keycloak at:", config.KeycloakURL)

	log.Println("Auth service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
