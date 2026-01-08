package e2e

import (
	"log"
	"os"
	"testing"

	"github.com/TATAROmangol/mess/e2e/config"
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

}
