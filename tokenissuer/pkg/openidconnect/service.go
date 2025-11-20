package openidconnect

import "context"

type Subject interface {
	GetID() string
	GetName() string
	GetEmail() string
}

type Service interface {
	VerifyToken(ctx context.Context, accessToken string) (Subject, error)
}
