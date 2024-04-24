// Package middleware http router middlewares
package middleware

import (
	"net/http"
)

// Auth middleware implementation
type Auth struct {
}

// NewAuth constructor
func NewAuth() *Auth {
	return &Auth{}
}

// Authenticate http middleware to verify wallet
func (a *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
