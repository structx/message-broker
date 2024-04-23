package domain

import (
	"context"

	"google.golang.org/grpc"
)

// AuthenticatorInterceptor
//
//go:generate mockery --name AuthenticatorInterceptor
type AuthenticatorInterceptor interface {
	// UnaryInterceptor
	UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	// StreamInterceptor
	StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}
