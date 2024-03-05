// Package domain application entities and interfaces
package domain

import "time"

// Block chain level block model
type Block struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Height    int64  `json:"height"`
	MaxHeight int64  `json:"max_height"`
	Timestamp string `json:"timestamp"`
	Txs       []*Tx  `json:"txs"`
}

// NewBlock return new block
func NewBlock(previousHash string, txs ...*Tx) *Block {
	return &Block{
		Hash:      "",
		PrevHash:  previousHash,
		Height:    int64(len(txs)),
		MaxHeight: 10,
		Timestamp: time.Now().Format(time.RFC3339),
		Txs:       txs,
	}
}
