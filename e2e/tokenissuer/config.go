package tokenissuer

import (
	"time"
)

type AdapterConfig struct {
	AuthURL       string        `yaml:"url"`
	ClientID      string        `yaml:"client_id"`
	ClientSecret  string        `yaml:"client_secret"`
	Login         string        `yaml:"login"`
	Password      string        `yaml:"password"`
	SubjectID     string        `yaml:"subject_id"`
	TokenDuration time.Duration `yaml:"token_duration"`
}

type IssuerConfig struct {
	VerifyGrpcAddress string `yaml:"verify_grpc_address"`
}

type Config struct {
	Adapter AdapterConfig `yaml:"adapter"`
	Issuer  IssuerConfig  `yaml:"issuer"`
}
