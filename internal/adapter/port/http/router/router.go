// Package router fuego server provider
package router

import (
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	pkgcontroller "github.com/trevatk/go-pkg/adapter/port/http/controller"

	"github.com/trevatk/mora/internal/adapter/port/http/controller"
	"github.com/trevatk/mora/internal/core/domain"
)

// NewRouter return new chi router
func NewRouter(logger *zap.Logger, raft domain.Raft) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	cc := []interface{}{
		pkgcontroller.NewBundle(logger),
		controller.NewRaft(logger, raft),
	}

	v1 := chi.NewRouter()

	for _, c := range cc {

		if c0, ok := c.(pkgcontroller.V0); ok {
			h := c0.RegisterRoutesV0()
			r.Mount("/", h)
		}

		if c1, ok := c.(pkgcontroller.V1); ok {
			h := c1.RegisterRoutesV1()
			v1.Mount("/", h)
		}

	}

	r.Mount("/api/v1", v1)

	return r
}
