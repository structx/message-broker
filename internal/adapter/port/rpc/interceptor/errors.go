package interceptor

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrMissingMetadata
	ErrMissingMetadata = status.Error(codes.InvalidArgument, "missing metadata")
)
