package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID = "app"
	clientSecret= "f3a485d3-8f82-482f-9844-ff7f269eab3c"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://local.dislackord.com:8080/auth/realms/dislackord")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		Endpoint: provider.Endpoint(),
		RedirectURL: "http://local.dislackord.com:8081/auth/callback",
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "magic"

	http.HandleFunc("/", func (resp http.ResponseWriter, req *http.Request) {
		http.Redirect(resp, req, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("state") != state {
			http.Error(resp, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, req.URL.Query().Get("code"))
		if err != nil {
			http.Error(resp, "failed to excahge token", http.StatusBadRequest)
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok  {
			http.Error(resp, "no id_token", http.StatusBadRequest)
			return
		}

		response := struct {
			OAuth2Token *oauth2.Token
			RawIDToken string
		}{
			oauth2Token, rawIDToken,
		}

		data, err := json.MarshalIndent(response, "","   ")
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		resp.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
