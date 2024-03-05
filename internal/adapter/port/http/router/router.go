// Package router fuego server provider
package router

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/go-fuego/fuego"

	"github.com/trevatk/block-broker/internal/adapter/port/http/controller"
	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/core/domain"
)

// NewRouter return new fuego server
func NewRouter(logger *zap.Logger, cfg *setup.Config, m domain.Messenger) *fuego.Server {

	s := fuego.NewServer(
		fuego.WithPort(fmt.Sprintf(":%s", cfg.Server.HTTPPort)),
		fuego.WithoutLogger(),
	)

	controllers := []interface{}{
		controller.NewMessages(logger, m),
		controller.NewHealthController(logger),
	}

	for _, c := range controllers {

		if v, ok := c.(controller.Controller); ok {
			v.RegisterRoutesV1(s)
		}

		if v, ok := c.(controller.ServiceController); ok {
			v.RegisterRoutesV0(s)
		}
	}

	return s
}
