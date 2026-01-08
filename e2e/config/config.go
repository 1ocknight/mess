package config

import (
	"flag"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
}

func LoadConfig() (*Config, error) {
	var configPath = flag.String("config-path", "", "path to config")
	flag.Parse()

	path := *configPath
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		panic("Config path is not set. Pass -config-path or set CONFIG_PATH")
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
