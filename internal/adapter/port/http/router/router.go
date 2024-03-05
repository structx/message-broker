// Package router fuego server provider
package router

import (
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"

	"github.com/trevatk/block-broker/internal/adapter/port/http/controller"
	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/core/domain"
)

// NewRouter return new fuego server
func NewRouter(logger *zap.Logger, cfg *setup.Config, m domain.Messenger) *echo.Echo {

	e := echo.New()

	controllers := []interface{}{
		controller.NewMessages(logger, m),
		controller.NewHealthController(logger),
	}

	g := e.Group("/")
	api := g.Group("/api")
	v1 := api.Group("/v1")

	for _, c := range controllers {

		if v, ok := c.(controller.Controller); ok {
			v.RegisterRoutesV1(v1)
		}

		if v, ok := c.(controller.ServiceController); ok {
			v.RegisterRoutesV0(g)
		}
	}

	return e
}
