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
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test"), 0)
		fmt.Println(h)
		prefix := strings.Repeat("0", 0)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("1", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test"), 1)
		fmt.Println(h)
		prefix := strings.Repeat("0", 1)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("2", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test"), 2)
		fmt.Println(h)
		prefix := strings.Repeat("0", 2)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("3", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test"), 3)
		fmt.Println(h)
		prefix := strings.Repeat("0", 3)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("4", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test"), 4)
		fmt.Println(h)
		prefix := strings.Repeat("0", 4)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
	t.Run("5", func(t *testing.T) {

		ts := time.Now().String()
		h := pow.HashWithSHA3AndDifficulty(ts, "ed64d4c5a4790be995081220a6fc59fea3b19d17551dc26bd3fcc116", []byte("test2"), 5)
		fmt.Println(h)
		prefix := strings.Repeat("0", 5)
		if !strings.HasPrefix(h, prefix) {
			t.Error("invalid hash generated")
		}
	})
}
