package jwksloader

import (
	"context"

	"github.com/TATAROmangol/mess/tokenissuer/pkg/jwks"
)

type Service interface {
	LoadJWKS(ctx context.Context) (map[string]jwks.JWKS, error)
}
