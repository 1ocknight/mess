package profile

import (
	"bytes"
	"testing"
	"time"

	"github.com/1ocknight/mess/e2e/utils"
	"github.com/1ocknight/mess/profile/pkg/dto"
)

type Tests struct {
	cfg *Config
}

func New(cfg *Config) *Tests {
	return &Tests{
		cfg: cfg,
	}
}

func (ts *Tests) cleanup(t *testing.T) {
	utils.GenericRequestWithAuth[*struct{}, dto.ProfileResponse](t, "DELETE", ts.cfg.DeleteProfileUrl, nil, utils.GetToken(t, ts.cfg.Auth))
}

func (ts *Tests) TestAddProfileWithUploadAvatar(t *testing.T) {
	defer ts.cleanup(t)

	alias := "test-alias"

	addReq := dto.AddProfileRequest{
		Alias: alias,
	}
	addResp := utils.GenericRequestWithAuth[dto.AddProfileRequest, dto.ProfileResponse](t, "POST", ts.cfg.AddProfileURL, &addReq, utils.GetToken(t, ts.cfg.Auth))
	if addResp == nil ||addResp.Alias != alias || addResp.Version != 1 || addResp.AvatarURL != "" {
		t.Fatalf("unexpected add profile response: %+v", addResp)
	}

	var content = []byte("test file content")

	file := bytes.NewReader(content)
	uploadUrlResp := utils.GenericRequestWithAuth[*struct{}, dto.UploadAvatarResponse](t, "PUT", ts.cfg.UploadAvatarURL, nil, utils.GetToken(t, ts.cfg.Auth))
	utils.SendFilePut(t, uploadUrlResp.UploadURL, file)

	time.Sleep(ts.cfg.WorkerDelay)
	getResp := utils.GenericRequestWithAuth[*struct{}, dto.ProfileResponse](t, "GET", ts.cfg.GetProfileURL, nil, utils.GetToken(t, ts.cfg.Auth))
	if getResp == nil ||getResp.Alias != alias || getResp.Version != 1 || getResp.AvatarURL == "" {
		t.Fatalf("unexpected get profile response: %+v", getResp)
	}

	utils.Checkfile(t, getResp.AvatarURL, content)
}

func (ts *Tests) Run(t *testing.T) {
	ts.cleanup(t)
	ts.TestAddProfileWithUploadAvatar(t)
}
