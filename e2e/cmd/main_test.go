package e2e

import (
	"log"
	"os"
	"testing"

	"github.com/TATAROmangol/mess/e2e/config"
	"github.com/TATAROmangol/mess/e2e/tokenissuer"
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
	ti := tokenissuer.New(&CFG.TokenIssuer)
	t.Run("token issuer", func(t *testing.T) {
		t.Parallel()
		ti.RunTests(t)
	})
}
