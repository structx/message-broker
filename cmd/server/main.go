// Package main entrypoint of application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/hashicorp/raft"
	"github.com/trevatk/go-pkg/logging"

	"github.com/trevatk/mora/internal/adapter/port/http/controller"
	"github.com/trevatk/mora/internal/adapter/port/http/middleware"
	"github.com/trevatk/mora/internal/adapter/port/http/router"
	"github.com/trevatk/mora/internal/adapter/port/http/server"
	"github.com/trevatk/mora/internal/adapter/port/rpc"
	"github.com/trevatk/mora/internal/adapter/port/rpc/interceptor"
	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/application"
	"github.com/trevatk/mora/internal/core/domain"
)

func main() {
	fx.New(
		fx.Provide(context.TODO),
		fx.Provide(setup.NewConfig),
		fx.Invoke(setup.DecodeHCLConfigFile),
		fx.Provide(logging.NewLoggerFromEnv),
		fx.Provide(fx.Annotate(application.NewRaftService, fx.As(new(domain.Raft)), fx.As(new(raft.FSM)))),
		fx.Provide(fx.Annotate(middleware.NewAuth, fx.As(new(domain.Authenticator)))),
		fx.Provide(fx.Annotate(interceptor.NewAuth, fx.As(new(domain.AuthenticatorInterceptor)))),
		fx.Provide(fx.Annotate(router.NewRouter, fx.As(new(http.Handler)))),
		fx.Provide(rpc.NewGRPCServer),
		fx.Provide(server.NewHTTPServer),
		fx.Invoke(controller.InvokeMetricsController),
		fx.Invoke(registerHooks),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}

func registerHooks(lc fx.Lifecycle, s1 *http.Server, s2 *rpc.GRPCServer) error {
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
				err := s2.Start()
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
