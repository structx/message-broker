package router_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/go-pkg/logging"
	"github.com/trevatk/mora/internal/adapter/port/http/router"
	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/domain"
)

func init() {
	_ = os.Setenv("LOG_PATH", "router.log")
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
}

func Test_NewRouter(t *testing.T) {
	t.Run("provider", func(t *testing.T) {

		assert := assert.New(t)

		ctx := context.TODO()

		logger, err := logging.NewLoggerFromEnv()
		assert.NoError(err)

		cfg := setup.NewConfig()
		assert.NoError(setup.ProcessConfigWithEnv(ctx, cfg))

		mockAuthenticator := domain.NewMockAuthenticator(t)

		s := router.NewRouter(logger, mockAuthenticator)
		assert.NotNil(s)
	})
}
