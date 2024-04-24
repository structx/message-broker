// Package interceptor gRPC request interceptors
package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// RecvMsg ...
func (w *wrappedStream) RecvMsg(m any) error {
	return w.ServerStream.RecvMsg(m)
}

// SendMessage ...
func (w *wrappedStream) SendMsg(m any) error {
	return w.ServerStream.SendMsg(m)
}

// Auth interceptor implementation
type Auth struct{}

// NewAuth constructor
func NewAuth() *Auth {
	return &Auth{}
}

// UnaryInterceptor single request interceptor to verify wallet access permissions
func (a *Auth) UnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingMetadata
	}

	m, err := handler(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "RPC failed with error %v", err)
	}

	return m, nil
}

// StreamInterceptor streaming request interceptor to verify wallet access permissions
func (a *Auth) StreamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	_, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return ErrMissingMetadata
	}

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		return status.Errorf(codes.Internal, "RPC failed with error %v", err)
	}

	return nil
}
