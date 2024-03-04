package kv

import "fmt"

// ErrNotFound key not found error
type ErrNotFound struct {
	Key []byte
}

// Error print error message
func (notFound *ErrNotFound) Error() string {
	return fmt.Sprintf("error key %s not found", string(notFound.Key))
}
