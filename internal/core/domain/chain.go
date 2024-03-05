package domain

// Chain blockchain interface
//
//go:generate mockery --name Chain
type Chain interface {
	// AddTransaction
	AddTx(data, action string, payload []byte, signature string) (string, error)
	// ListTransactions
	ListTransactions(limit, offset int) ([]*Tx, error)
	// ListTransactionsByAction
	ListTransactionsByAction(input string, limit, offset int) ([]*Tx, error)
}
