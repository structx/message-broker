// Package application logic
package application

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/hashicorp/raft"

	"github.com/trevatk/mora/internal/adapter/port/raftfx"
	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/domain"
	"github.com/trevatk/mora/pkg/messagebroker"
)

// RaftService implementation of raft interface
type RaftService struct {
	r *raft.Raft
}

// interface compliance
var _ domain.Raft = (*RaftService)(nil)
var _ raft.FSM = (*RaftService)(nil)

// NewRaftService constructor
func NewRaftService(cfg *setup.Config) (*RaftService, error) {

	rs := &RaftService{}

	r, _, err := raftfx.New(cfg, rs)
	if err != nil {
		return nil, fmt.Errorf("failed to create new raft %v", err)
	}
	rs.r = r

	return rs, nil
}

// Join raft
func (rs *RaftService) Join(_ context.Context, nm *domain.NewMember) (*domain.Member, error) {

	sID := raft.ServerID(nm.ServerID)
	sAddr := raft.ServerAddress(nm.Addr)

	prevIndex := rs.r.GetConfiguration().Index()
	f := rs.r.AddVoter(
		sID,
		sAddr,
		prevIndex,
		time.Second*15,
	)
	err := f.Error()
	if err != nil {
		return nil, fmt.Errorf("failed to add noter %v", err)
	}

	idx := f.Index()

	return &domain.Member{
		Index: idx,
	}, nil
}

// Notify all nodes in consensus of new message
func (rs *RaftService) Notify(_ context.Context, msg messagebroker.Msg) error {

	msgbytes, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("to marshal message %v", err)
	}

	f := rs.r.Apply(msgbytes, time.Second)
	err = f.Error()
	if err != nil {
		return rafterrors.MarkRetriable(err)
	}
	return nil
}

// Apply todo
func (rs *RaftService) Apply(_ *raft.Log) interface{} {
	// TODO: implement function
	return nil
}

// Snapshot todo
func (rs *RaftService) Snapshot() (raft.FSMSnapshot, error) {
	// TODO: implement function
	return nil, nil
}

// Restore todo
func (rs *RaftService) Restore(_ io.ReadCloser) error {
	// TODO: implement function
	return nil
}

// GetState getter current state of raft
func (rs *RaftService) GetState() domain.RaftState {
	s := rs.r.State()
	switch s {
	case raft.Leader:
		return domain.Leader
	default:
		return domain.Follower
	}
}
