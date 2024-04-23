package domain

import (
	"context"

	"github.com/trevatk/mora/pkg/messagebroker"
)

type RaftState int

const (
	Leader RaftState = iota
	Follower
	Voter
)

// NewMemember
type NewMember struct {
	ServerID string
	Addr     string
}

// Member
type Member struct {
	Index    uint64
	ServerID string
	Addr     string
}

// Raft
//
//go:generate mockery --name Raft
type Raft interface {
	Join(context.Context, *NewMember) (*Member, error)
	Apply(ctx context.Context, msg messagebroker.Msg) error
	GetState() RaftState
}
