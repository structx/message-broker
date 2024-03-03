package logging_test

import (
	"testing"

	"github.com/trevatk/block-broker/internal/adapter/logging"
)

func Test_NewLoggerFromEnv(t *testing.T) {
	t.Run("provider", func(t *testing.T) {
		logger, err := logging.NewLogger()
		if err != nil {
			t.Fatalf("failed to initialize new logger %v", err)
		}

		logger.Info("success")
	})
}
