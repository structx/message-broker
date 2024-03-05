// Package server http server provider
package server

import (
	"net/http"
	"time"

	"github.com/trevatk/block-broker/internal/adapter/setup"
)

// NewHTTPServer return new http server
func NewHTTPServer(cfg *setup.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
}
