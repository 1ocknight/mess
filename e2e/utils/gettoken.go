package utils

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

type GetTokensConfig struct {
	AuthURL      string `yaml:"url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Login        string `yaml:"login"`
	Password     string `yaml:"password"`
	SubjectID    string `yaml:"subject_id"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetToken(t *testing.T, cfg GetTokensConfig) string {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	client := resty.New()
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "password",
			"client_id":     cfg.ClientID,
			"client_secret": cfg.ClientSecret,
			"username":      cfg.Login,
			"password":      cfg.Password,
		}).
		Post(cfg.AuthURL)
	if err != nil {
		t.Fatalf("failed to request token: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var res TokenResponse
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		t.Fatalf("failed to unmarshal token response: %v", err)
	}

	return res.AccessToken
}
