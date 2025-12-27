package service

import (
	"context"

	"github.com/TATAROmangol/mess/profile/internal/model"
)

type Service interface {
	GetCurrentProfile(ctx context.Context) (*model.Profile, error)
	GetProfileFromSubjectID(ctx context.Context, subjID string) (*model.Profile, error)
	GetProfilesFromAlias(ctx context.Context, alias string, size int, token string) (string, []*model.Profile, error)
	GetAllProfiles(ctx context.Context, size int, token string) (string, []*model.Profile, error)

	AddProfile(ctx context.Context, alias string, avatar []byte, avatartType string) (*model.Profile, error)

	UpdateProfile(ctx context.Context, alias string) (*model.Profile, error)
	UpdateAvatar(ctx context.Context, avatar []byte) (string, error)

	DeleteCurrentProfile(ctx context.Context) error
	DeleteProfileFromSubjectID(ctx context.Context, subjID string) error
}
