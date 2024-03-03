package setup

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config
type Config struct {
	Server *Server
	KV     *KV
}

// NewConfig
func NewConfig() *Config {
	return &Config{
		Server: &Server{},
		KV:     &KV{},
	}
}

// ProcessConfigWithEnv
func ProcessConfigWithEnv(ctx context.Context, cfg *Config) error {
	return envconfig.Process(ctx, cfg)
}
