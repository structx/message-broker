package chain

import (
	"encoding/json"
	"fmt"

	"github.com/trevatk/block-broker/internal/core/domain"
	"github.com/trevatk/go-pkg/storage/kv"
)

type iterator struct {
	kv       kv.KV
	lastHash string
}

func newIterator(kv kv.KV, hash string) *iterator {
	return &iterator{
		kv:       kv,
		lastHash: hash,
	}
}

func (i *iterator) next() (*domain.Block, error) {

	if i.lastHash == "" {
		return nil, errEOC
	}

	value, err := i.kv.Get([]byte(i.lastHash))
	if err != nil {
		return nil, fmt.Errorf("unable to get block by hash %v", err)
	}

	var b domain.Block
	err = json.Unmarshal(value, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block %v", err)
	}

	i.lastHash = b.PrevHash

	return &b, nil
}
