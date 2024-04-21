package setup_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/mora/internal/adapter/setup"
)

func init() {
	_ = os.Setenv("KV_DIR", "./testfiles/kv")
}

func Test_ProcessConfigFromEnv(t *testing.T) {
	t.Run("process", func(t *testing.T) {

		assert := assert.New(t)

		cfg := setup.NewConfig()
		err := setup.ProcessConfigWithEnv(context.TODO(), cfg)
		assert.NoError(err)

		assert.Equal("./testfiles/kv", cfg.KV.Dir)
	})
}
