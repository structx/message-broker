// Package server http server provider
package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/trevatk/mora/internal/adapter/setup"
)

// NewHTTPServer return new http server
func NewHTTPServer(cfg *setup.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.BindAddr, fmt.Sprintf("%d", cfg.Server.Ports.HTTP)),
		Handler:      handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
}
