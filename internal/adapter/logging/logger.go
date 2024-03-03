// Package logging zap logger
package logging

import "go.uber.org/zap"

// NewLogger return new zap logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
