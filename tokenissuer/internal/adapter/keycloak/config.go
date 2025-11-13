package keycloak

import "time"

type Config struct {
	RefreshEndpoint      string
	ExchangeCodeEndpoint string
	ClientID             string
	ClientSecret         string
	Timeout              time.Duration
}
