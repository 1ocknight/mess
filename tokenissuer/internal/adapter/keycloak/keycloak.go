package keycloak

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	GrantType    = "authorization_code"
	RefreshToken = "refresh_token"

	RefreshTokenField = "refresh_token"
	GrantTypeField    = "grant_type"
	CodeField         = "code"
	RedirectURIField  = "redirect_uri"
	ClientIDField     = "client_id"
	ClientSecretField = "client_secret"
)

type Keycloak struct {
	cfg    Config
	client *resty.Client
}

func NewKeycloak(cfg Config) *Keycloak {
	client := resty.New()
	client.SetTimeout(cfg.Timeout)

	return &Keycloak{
		cfg:    cfg,
		client: client,
	}
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}

func (k *Keycloak) ExchangeCode(code string, redirectURL string) (*TokenResponse, error) {
	resp, err := k.client.R().
		SetFormData(map[string]string{
			GrantTypeField:    GrantType,
			CodeField:         code,
			RedirectURIField:  redirectURL,
			ClientIDField:     k.cfg.ClientID,
			ClientSecretField: k.cfg.ClientSecret,
		}).
		SetResult(&TokenResponse{}).
		Post(k.cfg.ExchangeCodeEndpoint)

	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("response error: %s", resp.String())
	}

	return resp.Result().(*TokenResponse), nil
}

func (k *Keycloak) Refresh(refreshToken string) (*TokenResponse, error) {
	resp, err := k.client.R().
		SetFormData(map[string]string{
			GrantTypeField:    RefreshToken,
			RefreshTokenField: refreshToken,
			ClientIDField:     k.cfg.ClientID,
			ClientSecretField: k.cfg.ClientSecret,
		}).
		SetResult(&TokenResponse{}).
		Post(k.cfg.RefreshEndpoint)

	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("response error: %s", resp.String())
	}

	return resp.Result().(*TokenResponse), nil
}
