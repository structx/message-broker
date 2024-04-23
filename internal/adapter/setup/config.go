// Package setup service configuration
package setup

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config service configuration
type Config struct {
	Server *Server
}

// NewConfig return new config class
func NewConfig() *Config {
	return &Config{
		Server: &Server{},
	}
}

// ProcessConfigWithEnv parse config from environment
func ProcessConfigWithEnv(ctx context.Context, cfg *Config) error {
	return envconfig.Process(ctx, cfg)
}
