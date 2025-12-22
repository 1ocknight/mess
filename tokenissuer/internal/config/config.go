package config

import (
	"fmt"
	"os"

	"github.com/TATAROmangol/mess/tokenissuer/internal/adapter/jwksloader/keycloak"
	"github.com/TATAROmangol/mess/tokenissuer/internal/service"
	"github.com/TATAROmangol/mess/tokenissuer/internal/transport/grpc"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Keycloak      keycloak.Config      `yaml:"keycloak"`
	GRPC          grpc.Config          `yaml:"grpc"`
	VerifyService service.VerifyConfig `yaml:"verify_config"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	return &cfg, nil
}
