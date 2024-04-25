package setup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/mora/internal/adapter/setup"
)

func init() {
	_ = os.Setenv("ROOT_CONFIG", "./testfiles/test_config.hcl")
}

func Test_ProcessConfigFromEnv(t *testing.T) {
	t.Run("process", func(t *testing.T) {

		assert := assert.New(t)

		cfg := setup.NewConfig()
		err := setup.DecodeHCLConfigFile(cfg)
		assert.NoError(err)

		fmt.Println(cfg.Server.BindAddr)
	})
}
