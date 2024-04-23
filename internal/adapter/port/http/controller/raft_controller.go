package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"

	"github.com/trevatk/go-pkg/http/controller"
	"github.com/trevatk/mora/internal/core/domain"
)

// Raft
type Raft struct {
	log *zap.SugaredLogger
	s   domain.Raft
}

// interface compliance
var _ controller.V1P = (*Raft)(nil)

// NewRaft
func NewRaft(logger *zap.Logger, raftService domain.Raft) *Raft {
	return &Raft{
		log: logger.Sugar().Named("RaftController"),
		s:   raftService,
	}
}

// RegisterRoutesV1P
func (rc *Raft) RegisterRoutesV1P() http.Handler {

	r := chi.NewRouter()

	r.Post("/raft/join", rc.Join)

	return r
}

// JoinParams
type JoinParams struct{}

// JoinResponse
type JoinResponse struct{}

// Join
func (rc *Raft) Join(w http.ResponseWriter, r *http.Request) {}
