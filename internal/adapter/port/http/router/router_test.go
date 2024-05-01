package router_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-pkg/adapter/logging"
	"github.com/structx/message-broker/internal/adapter/port/http/router"
	"github.com/structx/message-broker/internal/core/domain"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/test_config.hcl")
}

func Test_NewRouter(t *testing.T) {
	t.Run("provider", func(t *testing.T) {

		assert := assert.New(t)

		logger, err := logging.New(nil)
		assert.NoError(err)

		mockRaft := domain.NewMockRaft(t)

		s := router.NewRouter(logger, mockRaft)
		assert.NotNil(s)
	})
}
