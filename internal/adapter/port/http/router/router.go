// Package router fuego server provider
package router

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	pkgcontroller "github.com/trevatk/go-pkg/http/controller"
	"github.com/trevatk/mora/internal/adapter/port/http/controller"
	"github.com/trevatk/mora/internal/core/domain"
)

// NewRouter return new fuego server
func NewRouter(logger *zap.Logger, auth domain.Authenticator, raft domain.Raft) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	cc := []interface{}{
		pkgcontroller.NewBundle(logger),
		controller.NewRaft(logger, raft),
		controller.Metrics{}, // move to pkg bundle controller
	}

	v1 := chi.NewRouter()
	v1p := chi.NewRouter()

	v1p.Use(auth.Authenticate)

	for _, c := range cc {

		if c0, ok := c.(pkgcontroller.V0); ok {
			h := c0.RegisterRoutesV0()
			r.Mount("/", h)
		}

		if c1, ok := c.(pkgcontroller.V1); ok {
			h := c1.RegisterRoutesV1()
			v1.Mount("/", h)
		}

		if c1p, ok := c.(pkgcontroller.V1); ok {
			h := c1p.RegisterRoutesV1()
			v1p.Mount("/", h)
		}
	}

	r.Mount("/api/v1", v1)
	r.Mount("/api/v1/protected", v1p)

	return r
}
