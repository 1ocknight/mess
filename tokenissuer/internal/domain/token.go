package domain

import (
	"context"
	"fmt"
	"tokenissuer/internal/model"
	"tokenissuer/pkg/openidconnect"
)

type TokenService interface {
}

type TokenVerifier interface {
	VerifyToken(ctx context.Context, accessToken string) (*model.User, error)
}

type TokenDomain struct {
	ocid openidconnect.Service
}

func NewTokenDomain(ocid openidconnect.Service) *TokenDomain {
	return &TokenDomain{
		ocid: ocid,
	}
}

func (td *TokenDomain) VerifyToken(ctx context.Context, accessToken string) (*model.User, error) {
	respSubj, err := td.ocid.VerifyToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	subj := model.User{
		ID:    respSubj.GetID(),
		Name:  respSubj.GetName(),
		Email: respSubj.GetEmail(),
	}

	return &subj, nil
}
