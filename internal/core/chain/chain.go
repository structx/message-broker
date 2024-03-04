// Package chain blockchain logic
package chain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/trevatk/block-broker/internal/adapter/storage/kv"
	"github.com/trevatk/block-broker/internal/core/chain/pow"
	"github.com/trevatk/block-broker/internal/core/domain"
)

// Chain implementation of blockchain interface
type Chain struct {
	kv       domain.KV
	lastHash string
}

// interface compliance
var _ domain.Chain = (*Chain)(nil)

// NewChain return new chain class
// provides an existing chain or create new
func NewChain(db domain.KV) (*Chain, error) {

	v, err := db.Get([]byte("last_hash"))
	if err != nil {

		var notFound *kv.ErrNotFound
		if errors.As(err, &notFound) {
			// existing chain does not exist
			// create genesis block
			g := genesisBlock()
			blockbytes, err := json.Marshal(g)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal genesis block %v", err)
			}

			err = db.Put([]byte(g.Hash), blockbytes)
			if err != nil {
				return nil, fmt.Errorf("failed to put genesis block %v", err)
			}

			err = db.Put([]byte("last_hash"), []byte(g.Hash))
			if err != nil {
				return nil, fmt.Errorf("failed to put last_hash %v", err)
			}

			return &Chain{
				kv:       db,
				lastHash: g.Hash,
			}, nil
		}

		return nil, fmt.Errorf("failed to get last hash %v", err)

	}

	return &Chain{
		kv:       db,
		lastHash: hex.EncodeToString(v),
	}, nil
}

// AddBlock insert block
func (c *Chain) AddBlock(b *domain.Block) (string, error) {

	b.PrevHash = c.lastHash

	blockbytes, err := json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("failed to marshal block %v", err)
	}

	err = c.kv.Put([]byte(b.Hash), blockbytes)
	if err != nil {
		return "", fmt.Errorf("failed to put block bytes %v", err)
	}

	return b.Hash, nil
}

// GetBlock read block by key
func (c *Chain) GetBlock(key []byte) (*domain.Block, error) {

	k := []byte(c.lastHash)

	if !bytes.Equal(key, nil) {
		k = key
	}

	blockbytes, err := c.kv.Get(k)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by key %v", err)
	}

	var b domain.Block
	err = json.Unmarshal(blockbytes, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block bytes %v", err)
	}

	return &b, nil
}

// InsertTx add transaction to chain
func (c *Chain) InsertTx(tx *domain.Tx) (*domain.Tx, error) {

	// TODO:
	// add check for height of block >= max height
	// if at max height
	// 		create new block

	blockbytes, err := c.kv.Get([]byte(c.lastHash))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by key %v", err)
	}

	var b domain.Block
	err = json.Unmarshal(blockbytes, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block bytes %v", err)
	}

	b.Txs = append(b.Txs, tx)

	blockbytes, err = json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block %v", err)
	}

	err = c.kv.Put([]byte(b.Hash), blockbytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put block by key value %v", err)
	}

	return tx, nil
}

func genesisBlock() *domain.Block {
	return &domain.Block{
		Hash:      pow.HashWithSHA3AndDifficulty(time.Now().String(), 0),
		PrevHash:  "",
		Timestamp: time.Now().Format(time.RFC3339),
		Txs:       []*domain.Tx{},
	}
}
