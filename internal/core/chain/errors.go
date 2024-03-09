package chain

import "errors"

var (
	// ErrNotFound entity with hash not found
	ErrNotFound = errors.New("hash not found")

	errEOC = errors.New("iterator end of chain")
)
