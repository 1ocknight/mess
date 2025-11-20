package keycloak

import (
	"fmt"
	"time"

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

type Config struct {
	RefreshEndpoint      string        `json:"refresh_endpoint"`
	ExchangeCodeEndpoint string        `json:"exchange_code_endpoint"`
	ResetTokenEndpoint   string        `json:"reset_token_endpoint"`
	ClientID             string        `json:"client_id"`
	ClientSecret         string        `json:"client_secret"`
	Timeout              time.Duration `json:"timeout"`
}

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
		return nil, fmt.Errorf("response: %s", resp.String())
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
		return nil, fmt.Errorf("response: %s", resp.String())
	}

	return resp.Result().(*TokenResponse), nil
}

func (k *Keycloak) Logout(refreshToken string) error {
	resp, err := k.client.R().
		SetFormData(map[string]string{
			RefreshTokenField: refreshToken,
			ClientIDField:     k.cfg.ClientID,
			ClientSecretField: k.cfg.ClientSecret,
		}).
		Post(k.cfg.ResetTokenEndpoint)

	if err != nil {
		return fmt.Errorf("post logout: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("response: %s", resp.String())
	}

	return nil
}
