package e2e

import (
	"log"
	"os"
	"testing"

	"github.com/1ocknight/mess/e2e/config"
	"github.com/1ocknight/mess/e2e/profile"
)

var CFG *config.Config

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	CFG = cfg

	os.Exit(m.Run())
}

func TestRun(t *testing.T) {
	prof := profile.New(&CFG.Profile)
	prof.Run(t)
}
