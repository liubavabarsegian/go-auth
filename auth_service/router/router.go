package router

import (
	"authService/login"
	"authService/register"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

func SetUpRouter(keycloakClient *gocloak.GoCloak) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register.Register(w, r, keycloakClient)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login.Login(w, r, keycloakClient)
	})

	return mux
}
