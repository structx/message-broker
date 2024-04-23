package middleware

import "net/http"

// Authentication
type Authentication struct {
}

// Authenticate
func (a *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
