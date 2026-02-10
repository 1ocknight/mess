package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/1ocknight/mess/shared/verify"
	"github.com/1ocknight/mess/websocket/internal/hub/general"
	"github.com/1ocknight/mess/websocket/internal/transport"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Verify     verify.Config           `yaml:"verify"`
	HTTP       transport.HTTPConfig    `yaml:"http"`
	GeneralHub general.Config          `yaml:"general_hub"`
	Handler    transport.HandlerConfig `yaml:"handler"`
}

func LoadConfig() (*Config, error) {
	var configPath = flag.String("config", "", "path to config")
	flag.Parse()

	path := *configPath
	if path == "" {
		path = os.Getenv("CONFIG")
	}

	if path == "" {
		panic("Config path is not set. Pass -config or set CONFIG")
	}

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
