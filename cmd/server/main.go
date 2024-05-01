// Package main entrypoint of application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/multierr"

	"github.com/hashicorp/raft"

	"github.com/structx/go-pkg/adapter/logging"
	"github.com/structx/go-pkg/adapter/setup"
	pkgdomain "github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/decode"

	"github.com/structx/message-broker/internal/adapter/port/http/router"
	"github.com/structx/message-broker/internal/adapter/port/http/server"
	"github.com/structx/message-broker/internal/adapter/port/rpc"
	"github.com/structx/message-broker/internal/core/domain"
	"github.com/structx/message-broker/internal/core/service"
)

func main() {
	fx.New(
		fx.Provide(fx.Annotate(setup.New, fx.As(new(pkgdomain.Config)))),
		fx.Invoke(decode.ConfigFromEnv),
		fx.Provide(logging.New),
		fx.Provide(fx.Annotate(service.NewRaftService, fx.As(new(domain.Raft)), fx.As(new(raft.FSM)))),
		fx.Provide(fx.Annotate(router.NewRouter, fx.As(new(http.Handler)))),
		fx.Provide(rpc.NewGRPCServer),
		fx.Provide(server.NewHTTPServer),
		fx.Invoke(registerHooks),
	).Run()
}

func registerHooks(lc fx.Lifecycle, s1 *http.Server, s2 *rpc.GRPCServer, raftService domain.Raft) error {
	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {

				// start http server
				go func() {
					if err := s1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.Fatalf("failed to start http server %v", err)
					}
				}()

				// start gRPC server
				p := raftService.GetStartParams()
				err := s2.Start(p)
				if err != nil {
					return fmt.Errorf("failed to start rpc server %v", err)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {

				var result error

				// shutdown http server
				err := s1.Shutdown(ctx)
				if err != nil {
					result = multierr.Append(result, fmt.Errorf("failed to shutdown http server %v", err))
				}

				// graceful shutdown gRPC server
				s2.Shutdown()

				return result
			},
		},
	)
	return nil
}
