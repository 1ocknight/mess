package profile

import "github.com/TATAROmangol/mess/e2e/utils"

type Config struct {
	Auth             utils.GetTokensConfig `yaml:"auth"`
	DeleteProfileUrl string                `yaml:"delete_profile_url"`
	GetProfileURL    string                `yaml:"get_profile_url"`
	AddProfileURL    string                `yaml:"add_profile_url"`
	UploadAvatarURL  string                `yaml:"upload_avatar_url"`
}
