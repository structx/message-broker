package middleware

import (
	"net/http"
)

// Auth
type Auth struct {
}

// NewAuth
func NewAuth() *Auth {
	return &Auth{}
}

// Authenticate
func (a *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
