package domain

import "net/http"

// Authenticator
type Authenticator interface {
	// Authenticate
	Authenticate(next http.Handler) http.Handler
}
