package openidconnect

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
)

type Config struct {
	IssuerURL string `json:"issuer_url"`
	ClientID  string `json:"client_id"`
}

type Claims struct {
	SubjectID string `json:"sub"`
	Username  string `json:"preferred_username"`
	Email     string `json:"email"`
}

type OpenIDConnect struct {
	cfg      Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewOpenIDConnect(ctx context.Context, cfg Config) (*OpenIDConnect, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &OpenIDConnect{
		cfg:      cfg,
		provider: provider,
		verifier: verifier,
	}, nil
}

func (oc *OpenIDConnect) VerifyToken(ctx context.Context, accessToken string) (*Claims, error) {
	tok, err := oc.verifier.Verify(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	var claims Claims
	if err := tok.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}
