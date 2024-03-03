package domain

// Chain
//
//go:generate mockery --name Chain
type Chain interface {
	AddBlock(b *Block) (string, error)
	GetBlock(key []byte) (*Block, error)
	InsertTx(tx *Tx) (*Tx, error)
}
