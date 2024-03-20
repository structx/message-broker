package domain

// Chain blockchain interface
//
//go:generate mockery --name Chain
type Chain interface {
	// AddTx add new transaction
	AddTx(data, action string, payload []byte, signature string) (string, error)
	// ReadTx read transaction by id hash
	ReadTx(hash string) (*Tx, error)
	// ListTransactions
	ListTransactions(limit, offset int) ([]*Tx, error)
	// ListTransactionsByAction
	ListTransactionsByAction(input string, limit, offset int) ([]*Tx, error)
	// Shutdown
	Shutdown() error
}
