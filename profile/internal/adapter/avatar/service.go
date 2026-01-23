package avatar

import "context"

type Service interface {
	GetUploadURL(ctx context.Context, subjID string) (string, error)
	GetAvatarURL(ctx context.Context, subjID string) (string, error)
	DeleteObjects(ctx context.Context, subjIDs []string) error
}
