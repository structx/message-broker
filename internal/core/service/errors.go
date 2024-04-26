package service

import "errors"

var (
	// ErrNotFound entity identifier not found
	ErrNotFound = errors.New("entity not found")
)
