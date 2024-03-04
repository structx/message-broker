package domain

// KV key value interface
type KV interface {
	// Get value by key
	Get(key []byte) ([]byte, error)
	// Put set key/value pair
	Put(key, value []byte) error
	// Close database connection
	Close() error
}
