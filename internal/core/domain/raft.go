// Package domain application models and interfaces
package domain

import (
	"context"

	"github.com/trevatk/mora/pkg/messagebroker"
)

// RaftState current raft state
type RaftState int

const (
	// Leader leader
	Leader RaftState = iota
	// Follower follower
	Follower
	// Voter voter
	Voter
)

// NewMember application model new node
type NewMember struct {
	ServerID string
	Addr     string
}

// Member application model existing node
type Member struct {
	Index    uint64
	ServerID string
	Addr     string
}

// Raft application logic raft interface
//
//go:generate mockery --name Raft
type Raft interface {
	// Join raft
	Join(context.Context, *NewMember) (*Member, error)
	// Notify nodes in consensus
	Notify(ctx context.Context, msg messagebroker.Msg) error
	// GetState getter raft state
	GetState() RaftState
}
