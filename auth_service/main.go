package main

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func main() {
	client := gocloak.NewClient("http://localhost:8080")
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, "admin", "admin", "admin_realm")
	if err != nil {
		fmt.Println(err)
		panic("Something wrong with the credentials or url")
	}

	user := gocloak.User{
		FirstName: gocloak.StringP("Bob"),
		LastName:  gocloak.StringP("Uncle"),
		Email:     gocloak.StringP("something@really.wrong"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("CoolGuy"),
	}

	_, err = client.CreateUser(ctx, token.AccessToken, "admin_realm", user)
	if err != nil {
		panic("Oh no!, failed to create user :(")
	}
}
