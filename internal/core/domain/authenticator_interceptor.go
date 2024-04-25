package domain

import (
	"context"

	"google.golang.org/grpc"
)

// AuthenticatorInterceptor gRPC auth interceptor interface
//
//go:generate mockery --name AuthenticatorInterceptor
type AuthenticatorInterceptor interface {
	// UnaryInterceptor single request wallet verification
	UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	// StreamInterceptor streaming wallet verification
	StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}
