package domain

import "net/http"

// Authenticator http auth middleware interface
//
//go:generate mockery --name Authenticator
type Authenticator interface {
	// Authenticate verify wallet access permissions
	Authenticate(next http.Handler) http.Handler
}
