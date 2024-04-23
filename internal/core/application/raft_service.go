package application

import (
	"github.com/hashicorp/raft"
	"github.com/trevatk/mora/internal/core/domain"
)

// RaftService
type RaftService struct {
	r *raft.Raft
}

// interface compliance
var _ domain.Raft = (*RaftService)(nil)
