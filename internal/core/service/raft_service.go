// Package service application logic
package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"

	"github.com/trevatk/go-pkg/adapter/port/raftfx"
	pkgdomain "github.com/trevatk/go-pkg/domain"
	"github.com/trevatk/mora/internal/core/domain"
)

// RaftService implementation of raft interface
type RaftService struct {
	r  *raft.Raft
	tm *transport.Manager
}

// interface compliance
var _ domain.Raft = (*RaftService)(nil)
var _ raft.FSM = (*RaftService)(nil)

// NewRaftService constructor
func NewRaftService(cfg pkgdomain.Config) (*RaftService, error) {

	rs := &RaftService{}

	r, tm, err := raftfx.New(cfg, rs)
	if err != nil {
		return nil, fmt.Errorf("failed to create new raft %v", err)
	}
	rs.r = r
	rs.tm = tm

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
func (rs *RaftService) Notify(_ context.Context, msg pkgdomain.Envelope) error {

	f := rs.r.Apply(msg.GetPayload(), time.Second)
	err := f.Error()
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

// GetStartParams getter gRPC start params
func (rs *RaftService) GetStartParams() *domain.GrpcStartParams {
	return &domain.GrpcStartParams{
		TransportManager: rs.tm,
		Raft:             rs.r,
	}
}
