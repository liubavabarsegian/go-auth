package login

import (
	"authService/config"
	"context"
	"log"
	"testing"

	"github.com/Nerzal/gocloak/v13"
)

func TestLogin(t *testing.T) {
	keycloakClient := gocloak.NewClient(config.KeycloakURL)
	keycloakClient.RestyClient().SetDebug(true)

	if keycloakClient == nil {
		t.Fatal("Keycloak client is nil")
	}

	t.Run("Successful_Login", func(t *testing.T) {
		// Пример: параметры для логина
		username := "auth-admin"
		password := "admin"

		token, err := keycloakClient.Login(
			context.Background(),
			config.ClientID,
			config.ClientSecret,
			config.Realm,
			username,
			password,
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if token.AccessToken == "" {
			t.Fatal("expected token to be non-empty")
		}

		log.Printf("Access Token: %s", token.AccessToken)
	})
}
