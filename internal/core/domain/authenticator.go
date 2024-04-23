package domain

import "net/http"

// Authenticator
//
//go:generate mockery --name Authenticator
type Authenticator interface {
	// Authenticate
	Authenticate(next http.Handler) http.Handler
}
