// Package server http server provider
package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	pkgdomain "github.com/structx/go-pkg/domain"
)

// NewHTTPServer return new http server
func NewHTTPServer(cfg pkgdomain.Config, handler http.Handler) *http.Server {
	scfg := cfg.GetServer()
	return &http.Server{
		Addr:         net.JoinHostPort(scfg.BindAddr, fmt.Sprintf("%d", scfg.Ports.HTTP)),
		Handler:      handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
}
