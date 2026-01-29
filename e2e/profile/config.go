package profile

import (
	"time"

	"github.com/1ocknight/mess/e2e/utils"
)

type Config struct {
	Auth             utils.GetTokensConfig `yaml:"auth"`
	WorkerDelay      time.Duration         `yaml:"worker_delay"`
	UpdateProfileURL string                `yaml:"update_profile_url"`
	DeleteProfileUrl string                `yaml:"delete_profile_url"`
	GetProfileURL    string                `yaml:"get_profile_url"`
	AddProfileURL    string                `yaml:"add_profile_url"`
	UploadAvatarURL  string                `yaml:"upload_avatar_url"`
}
