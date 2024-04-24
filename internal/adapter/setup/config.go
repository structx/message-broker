// Package setup service configuration
package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	hcl2 "github.com/hashicorp/hcl/v2/hclsimple"
)

// Config service configuration
type Config struct {
	Server Server `hcl:"server,block"`
	Raft   Raft   `hcl:"raft,block"`
	Logger Logger `hcl:"logger,block"`
}

// NewConfig constructor
func NewConfig() *Config {
	return &Config{
		Server: Server{},
		Raft:   Raft{},
	}
}

// DecodeHCLConfigFile read config path from environment and parse file
func DecodeHCLConfigFile(cfg *Config) error {

	configFile := os.Getenv("DSERVICE_CONFIG")
	if configFile == "" {
		return errors.New("$DSERVICE_CONFIG must be set")
	}

	if err := hcl2.DecodeFile(
		filepath.Clean(configFile),
		nil,
		cfg,
	); err != nil {
		return fmt.Errorf("failed decode config file %v", err)
	}

	return nil
}
