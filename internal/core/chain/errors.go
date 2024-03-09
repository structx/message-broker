package chain

import "errors"

var (
	// ErrNotFound
	ErrNotFound = errors.New("hash not found")

	errEOC = errors.New("iterator end of chain")
)
