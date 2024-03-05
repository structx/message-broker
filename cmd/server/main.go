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

	"github.com/trevatk/block-broker/internal/adapter/logging"
	"github.com/trevatk/block-broker/internal/adapter/port/http/router"
	"github.com/trevatk/block-broker/internal/adapter/port/http/server"
	"github.com/trevatk/block-broker/internal/adapter/port/rpc"
	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/adapter/storage/kv"
	"github.com/trevatk/block-broker/internal/core/application"
	"github.com/trevatk/block-broker/internal/core/chain"
	"github.com/trevatk/block-broker/internal/core/domain"
)

func main() {
	fx.New(
		fx.Provide(context.TODO),
		fx.Provide(logging.NewLogger),
		fx.Provide(setup.NewConfig),
		fx.Invoke(setup.ProcessConfigWithEnv),
		fx.Provide(fx.Annotate(kv.NewPebble, fx.As(new(domain.KV)))),
		fx.Provide(fx.Annotate(chain.NewChain, fx.As(new(domain.Chain)))),
		fx.Provide(fx.Annotate(application.NewMessagingService, fx.As(new(domain.Messenger)))),
		fx.Provide(fx.Annotate(router.NewRouter, fx.As(new(http.Handler)))),
		fx.Provide(rpc.NewGRPCServer),
		fx.Provide(server.NewHTTPServer),
		fx.Invoke(registerHooks),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}

func registerHooks(lc fx.Lifecycle, s1 *http.Server, s2 *rpc.GRPCServer, kv domain.KV) error {
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

				// close kv database
				err = kv.Close()
				if err != nil {
					result = multierr.Append(result, fmt.Errorf("failed to close kv database %v", err))
				}

				return result
			},
		},
	)
	return nil
}
