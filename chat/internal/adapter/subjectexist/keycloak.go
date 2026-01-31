package subjectexist

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	KeycloakURL  string        `yaml:"keycloak_url"`
	Realm        string        `yaml:"realm"`
	ClientID     string        `yaml:"client_id"`
	ClientSecret string        `yaml:"client_secret"`
	Timeout      time.Duration `yaml:"timeout"`
}

type Keycloak struct {
	cfg    Config
	client *resty.Client
	oauth  *clientcredentials.Config
}

func New(cfg Config) (*Keycloak, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		cfg.KeycloakURL, cfg.Realm)

	oauthConfig := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL,
	}

	client := resty.New()
	client.SetTimeout(cfg.Timeout)

	return &Keycloak{
		cfg:    cfg,
		client: client,
		oauth:  oauthConfig,
	}, nil
}

func (k *Keycloak) SubjectExists(ctx context.Context, subjectID string) (bool, error) {
	token, err := k.oauth.Token(ctx)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s",
		k.cfg.KeycloakURL, k.cfg.Realm, subjectID)

	resp, err := k.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		Get(url)

	if err != nil {
		return false, err
	}

	switch resp.StatusCode() {
	case 200:
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status %d: %s",
			resp.StatusCode(), resp.Body())
	}
}
