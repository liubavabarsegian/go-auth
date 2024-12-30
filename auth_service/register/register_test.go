package register

import (
	"authService/config"
	"context"
	"log"
	"testing"

	"github.com/Nerzal/gocloak/v13"
)

func TestRegister(t *testing.T) {
	// Инициализация keycloakClient
	keycloakClient := gocloak.NewClient(config.KeycloakURL)
	keycloakClient.RestyClient().SetDebug(true)

	if keycloakClient == nil {
		t.Fatal("Keycloak client is nil")
	}

	t.Run("Successful_Registration", func(t *testing.T) {
		// Пример данных для регистрации
		username := "newuser"
		email := "newuser@example.com"
		password := "password123"

		// Получение токена администратора для выполнения регистрации
		adminToken, err := keycloakClient.LoginAdmin(
			context.Background(),
			"admin",  // Имя администратора
			"admin",  // Пароль администратора
			"master", // Realm администратора
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Данные пользователя для регистрации
		user := gocloak.User{
			Username: &username,
			Email:    &email,
			Enabled:  gocloak.BoolP(true),
			Credentials: &[]gocloak.CredentialRepresentation{
				{
					Type:      gocloak.StringP("password"),
					Value:     &password,
					Temporary: gocloak.BoolP(false),
				},
			},
		}

		// Создание пользователя в Keycloak
		_, err = keycloakClient.CreateUser(context.Background(), adminToken.AccessToken, config.Realm, user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Логирование успешного результата
		log.Println("User successfully registered")
	})
}
