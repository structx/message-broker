// Package kv database implementation
package kv

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	pdb "github.com/cockroachdb/pebble"

	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/core/domain"
)

// Pebble db wrapper class
type Pebble struct {
	db *pdb.DB
}

// interface compliance
var _ domain.KV = (*Pebble)(nil)

// NewPebble return new pebble db wrapper class
func NewPebble(cfg *setup.Config) (*Pebble, error) {

	opts := &pdb.Options{}
	db, err := pdb.Open(cfg.KV.Dir, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open pebble db: %v", err)
	}

	return &Pebble{
		db: db,
	}, nil
}

// Put set key/value pair
func (p *Pebble) Put(key, value []byte) error {
	return p.db.Set(key, value, pebble.Sync)
}

// Get value by key
func (p *Pebble) Get(key []byte) ([]byte, error) {

	v, closer, err := p.db.Get(key)
	if err != nil && err == pdb.ErrNotFound {
		return []byte{}, &ErrNotFound{Key: key}
	} else if err != nil {
		return []byte{}, fmt.Errorf("failed to get key value %v", err)
	}

	err = closer.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to close closer %v", err)
	}

	return v, nil
}

// Close database connection
func (p *Pebble) Close() error {
	return p.db.Close()
}
