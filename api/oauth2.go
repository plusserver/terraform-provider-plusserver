package api

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	TokenURL     string
}

func NewClient(config *OAuthConfig) (*http.Client, error) {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL:  config.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.PasswordCredentialsToken(ctx, config.Username, config.Password)
	if err != nil {
		return nil, err
	}
	client := conf.Client(ctx, tok)
	return client, nil
}
