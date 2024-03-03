// package kv
package kv

import "fmt"

// ErrNotFound
type ErrNotFound struct {
	Key []byte
}

// Error
func (notFound *ErrNotFound) Error() string {
	return fmt.Sprintf("error key %s not found", string(notFound.Key))
}
