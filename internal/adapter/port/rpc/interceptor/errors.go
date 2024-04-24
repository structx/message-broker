package interceptor

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrMissingMetadata invalid or missing metadata
	ErrMissingMetadata = status.Error(codes.InvalidArgument, "missing metadata")
)
