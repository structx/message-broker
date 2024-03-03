package pow_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/trevatk/block-broker/internal/core/chain/pow"
)

func TestHashWithSHA3AndDifficulty(t *testing.T) {
	t.Run("0", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, 0)
		fmt.Println(h)
		prefix := strings.Repeat("0", 0)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("1", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, 1)
		fmt.Println(h)
		prefix := strings.Repeat("0", 1)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("2", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, 2)
		fmt.Println(h)
		prefix := strings.Repeat("0", 2)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("3", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, 3)
		fmt.Println(h)
		prefix := strings.Repeat("0", 3)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
}
