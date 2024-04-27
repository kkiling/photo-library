package session_manager

import "time"

// SslConfig .
type SslConfig struct {
	CertFile       string `yaml:"cert_file"`
	PrivateKeyFile string `yaml:"private_key_file"`
	PublicKeyFile  string `yaml:"public_key_file"`
}

// SessionConfig конфигурация выпуска jwt
type SessionConfig struct {
	Audience             string        `yaml:"audience"`
	Issuer               string        `yaml:"issuer"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}
