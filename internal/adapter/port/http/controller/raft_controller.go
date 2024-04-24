// Package controller exposed http endpoints and models
package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"

	"github.com/trevatk/go-pkg/http/controller"
	"github.com/trevatk/mora/internal/core/domain"
)

// Raft controller
type Raft struct {
	log *zap.SugaredLogger
	s   domain.Raft
}

// interface compliance
var _ controller.V1P = (*Raft)(nil)

// NewRaft constructor
func NewRaft(logger *zap.Logger, raftService domain.Raft) *Raft {
	return &Raft{
		log: logger.Sugar().Named("RaftController"),
		s:   raftService,
	}
}

// RegisterRoutesV1P return router for protected endpoints in raft controller
func (rc *Raft) RegisterRoutesV1P() http.Handler {

	r := chi.NewRouter()

	r.Post("/raft/join", rc.Join)

	return r
}

// JoinRequestParams exposed http model for join handler
type JoinRequestParams struct{}

// JoinResponse exposed http response model for join handler
type JoinResponse struct{}

// Join raft handler
func (rc *Raft) Join(_ http.ResponseWriter, _ *http.Request) {
	// TODO: implement handler
}
