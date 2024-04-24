package router_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/go-pkg/logging"
	"github.com/trevatk/mora/internal/adapter/port/http/router"
	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/domain"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/test_config.hcl")
}

func Test_NewRouter(t *testing.T) {
	t.Run("provider", func(t *testing.T) {

		assert := assert.New(t)

		logger, err := logging.NewLoggerFromEnv()
		assert.NoError(err)

		cfg := setup.NewConfig()
		assert.NoError(setup.DecodeHCLConfigFile(cfg))

		mockAuthenticator := domain.NewMockAuthenticator(t)
		mockRaft := domain.NewMockRaft(t)

		s := router.NewRouter(logger, mockAuthenticator, mockRaft)
		assert.NotNil(s)
	})
}
