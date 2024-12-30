package config

import (
	"os"
)

var (
	Realm        = os.Getenv("KEYCLOAK_REALM")
	ClientID     = os.Getenv("KEYCLOAK_CLIENT_ID")
	ClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	KeycloakURL  = os.Getenv("KEYCLOAK_URL")
)
