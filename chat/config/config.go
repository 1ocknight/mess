package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/1ocknight/mess/chat/internal/adapter/subjectexist"
	"github.com/1ocknight/mess/chat/internal/transport"
	messagesenderworker "github.com/1ocknight/mess/chat/internal/worker/messagesender"
	"github.com/1ocknight/mess/shared/postgres"
	"github.com/1ocknight/mess/shared/redisclient"
	"github.com/1ocknight/mess/shared/verify"
	"github.com/goccy/go-yaml"
)

type Config struct {
	MigrationsPath      string                     `yaml:"migrations_path"`
	Postgres            postgres.Config            `yaml:"postgres"`
	HTTP                transport.Config           `yaml:"http"`
	LoggerDebug         bool                       `yaml:"logger_debug"`
	LastReadSender      redisclient.Config         `yaml:"last_read_sender"`
	MessageSender       redisclient.Config         `yaml:"message_sender"`
	SubjectExist        subjectexist.Config        `yaml:"subject_exist"`
	MessageSenderWorker messagesenderworker.Config `yaml:"message_sender_worker"`
	Verify              verify.Config              `yaml:"verify"`
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
