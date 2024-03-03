package domain

type KV interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Close() error
}
