package service

import (
	"context"
	"fmt"

	"github.com/TATAROmangol/mess/profile/internal/ctxkey"
	"github.com/TATAROmangol/mess/profile/internal/model"
	"github.com/TATAROmangol/mess/profile/internal/storage"
	"github.com/TATAROmangol/mess/profile/internal/storage/profile"
)

var (
	DefaultPageSize = 100

	Asc       = true
	SortLabel = profile.AliasLabel
)

type IMPL struct {
	s storage.Service
}

func New(s storage.Service) *IMPL {
	return &IMPL{
		s: s,
	}
}

func (i *IMPL) GetCurrentProfile(ctx context.Context) (*model.Profile, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %v", err)
	}

	return i.GetProfileFromSubjectID(ctx, subj.GetSubjectId())
}

func (i *IMPL) GetProfileFromSubjectID(ctx context.Context, subjID string) (*model.Profile, error) {
	profile, err := i.s.Profile().GetProfileFromSubjectID(ctx, subjID)
	if err != nil {
		return nil, fmt.Errorf("get profile from subject id: %v", err)
	}

	return profile, nil
}

func (i *IMPL) GetProfilesFromAlias(ctx context.Context, alias string, size int, token string) (string, []*model.Profile, error) {
	if token != "" {
		token, profiles, err := i.s.Profile().GetProfilesFromAliasWithToken(ctx, alias, token)
		if err != nil {
			return "", nil, fmt.Errorf("get profiles from alias with token: %v", err)
		}
		return token, profiles, nil
	}

	if size == 0 {
		size = DefaultPageSize
	}

	token, profiles, err := i.s.Profile().GetProfilesFromAlias(ctx, size, Asc, SortLabel, alias)
	if err != nil {
		return "", nil, fmt.Errorf("first get profiles from alias: %v", err)
	}

	return token, profiles, nil
}

func (i *IMPL) GetAllProfiles(ctx context.Context, size int, token string) (string, []*model.Profile, error) {
	token, profiles, err := i.GetProfilesFromAlias(ctx, "", size, token)
	if err != nil {
		return "", nil, fmt.Errorf("get profiles from alias: %v", err)
	}

	return token, profiles, nil
}

func (i *IMPL) AddProfile(ctx context.Context, alias string, avatar []byte, avatartType string) (*model.Profile, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %v", err)
	}

	avatarUrl, err := i.s.Avatar().Upload(ctx, subj.GetSubjectId(), avatar, avatartType)
	if err != nil {
		return nil, fmt.Errorf("avatar upload: %v", err)
	}

	profile, err := i.s.Profile().AddProfile(ctx, subj.GetSubjectId(), alias, avatarUrl)
	if err != nil {
		return nil, fmt.Errorf("add profile: %v", err)
	}

	return profile, nil
}
