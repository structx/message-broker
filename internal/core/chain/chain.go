// Package chain blockchain logic
package chain

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/trevatk/block-broker/internal/core/chain/pow"
	"github.com/trevatk/block-broker/internal/core/domain"
	"github.com/trevatk/go-pkg/storage/kv"
)

// Chain implementation of blockchain interface
type Chain struct {
	kv       kv.KV
	lastHash string
}

// interface compliance
var _ domain.Chain = (*Chain)(nil)

// NewChain return new chain class
// provides an existing chain or create new
func NewChain(db kv.KV) (*Chain, error) {

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

// AddTx insert new transaction
func (c *Chain) AddTx(data, action string, payload []byte, signature string) (string, error) {

	// create new transactions
	tx := domain.NewTx(data, action, payload, signature)
	tx.SetID()

	// read current working block
	currentblock, err := c.getCurrentBlock()
	if err != nil {
		return "", fmt.Errorf("failed to get current block %v", err)
	}

	// check if working block has reached max height
	if currentblock.Height >= currentblock.MaxHeight {

		// create new block
		currentblock = domain.NewBlock(c.lastHash, tx)
		// generate hash with simple difficulty
		hash := pow.HashWithSHA3AndDifficulty(currentblock.Timestamp, currentblock.PrevHash, tx.Payload, 1)

		// assign hash and chain latest hash
		currentblock.Hash = hash
		c.lastHash = currentblock.Hash
	}

	currentblock.Txs = append(currentblock.Txs, tx)

	value, err := json.Marshal(currentblock)
	if err != nil {
		return "", fmt.Errorf("failed to marshal current block %v", err)
	}

	err = c.kv.Put([]byte(currentblock.Hash), value)
	if err != nil {
		return "", fmt.Errorf("failed to set current block value %v", err)
	}

	return hex.EncodeToString(tx.ID), nil
}

// ReadTx read transaction by id hash
func (c *Chain) ReadTx(hash string) (*domain.Tx, error) {

	it := newIterator(c.kv, c.lastHash)

	for {

		block, err := it.next()
		if err != nil && err == errEOC {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("failed to read next transaction %v", err)
		}

		if block == nil {
			break
		}

	INNER:
		for _, tx := range block.Txs {

			if hex.EncodeToString(tx.ID) == hash {
				return tx, nil
			}

			continue INNER
		}
	}

	return nil, ErrNotFound
}

// ListTransactions with limit and offset
func (c *Chain) ListTransactions(limit, offset int) ([]*domain.Tx, error) {

	it := newIterator(c.kv, c.lastHash)

	txSlice := make([]*domain.Tx, 0, limit)

	counter := 0

OUTER:
	for {

		block, err := it.next()
		if err != nil && err == errEOC {
			break OUTER
		} else if err != nil {
			return nil, fmt.Errorf("failed to read next transaction %v", err)
		}

		if block == nil {
			break
		}

	INNER:
		for _, tx := range block.Txs {

			if offset >= 0 {
				offset--
				continue INNER
			}

			if counter == limit {
				break OUTER
			}

			txSlice = append(txSlice, tx)
			counter++
		}
	}

	return txSlice, nil
}

// ListTransactionsByAction read all transactions with matching actions
func (c *Chain) ListTransactionsByAction(input string, limit, offset int) ([]*domain.Tx, error) {

	it := newIterator(c.kv, c.lastHash)

	txSlice := make([]*domain.Tx, 0, limit)

	counter := 0

OUTER:
	for {

		block, err := it.next()
		if err != nil && err == errEOC {
			break OUTER
		} else if err != nil {
			return nil, fmt.Errorf("failed to read next transaction %v", err)
		}

		if block == nil {
			break
		}

	INNER:
		for _, tx := range block.Txs {

			if offset >= 0 {
				offset--
				continue INNER
			}

			if counter == limit {
				break OUTER
			}

			if tx.Action == input {
				txSlice = append(txSlice, tx)
				counter++
			}
		}
	}

	return txSlice, nil
}

func (c *Chain) getCurrentBlock() (*domain.Block, error) {

	value, err := c.kv.Get([]byte(c.lastHash))
	if err != nil {
		return nil, fmt.Errorf("failed to get last block with hash %v", err)
	}

	var b domain.Block
	err = json.Unmarshal(value, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block %v", err)
	}

	return &b, nil
}

func genesisBlock() *domain.Block {
	return &domain.Block{
		Hash:      pow.HashWithSHA3AndDifficulty(time.Now().UTC().String(), "", []byte("gensis block"), 0),
		PrevHash:  "",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Txs:       []*domain.Tx{},
		Height:    1,
		MaxHeight: 10,
	}
}
