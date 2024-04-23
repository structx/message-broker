package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/hashicorp/raft"
	"github.com/trevatk/mora/internal/core/domain"
	"github.com/trevatk/mora/pkg/messagebroker"
)

// RaftService
type RaftService struct {
	r *raft.Raft
}

// interface compliance
var _ domain.Raft = (*RaftService)(nil)

// NewRaftService
func NewRaftService(raft *raft.Raft) *RaftService {
	return &RaftService{
		r: raft,
	}
}

// Join
func (rs *RaftService) Join(ctx context.Context, nm *domain.NewMember) (*domain.Member, error) {

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

// Apply
func (rs *RaftService) Apply(ctx context.Context, msg messagebroker.Msg) error {

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

// GetState
func (rs *RaftService) GetState() domain.RaftState {
	s := rs.r.State()
	switch s {
	case raft.Leader:
		return domain.Leader
	default:
		return domain.Follower
	}
}
