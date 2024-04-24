package domain

import (
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
)

// GrpcStartParams gRPC server start params
type GrpcStartParams struct {
	TransportManager *transport.Manager
	Raft             *raft.Raft
}
