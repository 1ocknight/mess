package keycloak_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/TATAROmangol/mess/shared/auth/keycloak"
	loggermocks "github.com/TATAROmangol/mess/shared/logger/mocks"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	AuthURL      string `yaml:"url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Login        string `yaml:"login"`
	Password     string `yaml:"password"`
	SubjectID    string `yaml:"subject_id"`
	JWKSEndpoint string `yaml:"jwks_endpoint"`
}

var CFG *TestConfig

func TestMain(m *testing.M) {
	CFG = &TestConfig{
		AuthURL:      "http://localhost:7070/realms/e2e-realm/protocol/openid-connect/token",
		ClientID:     "e2e",
		ClientSecret: "e2e",
		Login:        "e2e",
		Password:     "e2e",
		SubjectID:    "8f446bd3-7c32-4bbe-aa0d-516d621cf208",
		JWKSEndpoint: "http://localhost:7070/realms/e2e-realm/protocol/openid-connect/certs",
	}
}

func getToken(t *testing.T) string {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	client := resty.New()

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "password",
			"client_id":     CFG.ClientID,
			"client_secret": CFG.ClientSecret,
			"username":      CFG.Login,
			"password":      CFG.Password,
		}).
		Post(CFG.AuthURL)

	if err != nil {
		t.Fatalf("failed to request token: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var res string

	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		t.Fatalf("failed to unmarshal token response: %v", err)
	}

	return res
}

func TestKeycloak_Verify(t *testing.T) {
	ctrl := gomock.NewController(t)
	lg := loggermocks.NewMockLogger(ctrl)

	k, err := keycloak.New(keycloak.Config{JWKSEndpoint: CFG.JWKSEndpoint}, lg)
	if err != nil {
		t.Fatalf("keycloak new: %v", err)
	}

	token := getToken(t)
	subj, err := k.Verify(token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}

	require.Equal(t, subj.GetSubjectId(), CFG.SubjectID)
}
