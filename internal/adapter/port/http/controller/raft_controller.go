// Package controller exposed http endpoints and models
package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"go.uber.org/zap"

	"github.com/structx/go-pkg/adapter/port/http/controller"
	"github.com/structx/message-broker/internal/core/domain"
)

// Raft controller
type Raft struct {
	log     *zap.SugaredLogger
	service domain.Raft
}

// interface compliance
var _ controller.V1P = (*Raft)(nil)

// NewRaft constructor
func NewRaft(logger *zap.Logger, raftService domain.Raft) *Raft {
	return &Raft{
		log:     logger.Sugar().Named("RaftController"),
		service: raftService,
	}
}

// RegisterRoutesV1P return router for protected endpoints in raft controller
func (rc *Raft) RegisterRoutesV1P() http.Handler {

	r := chi.NewRouter()

	r.Put("/raft/join", rc.Join)

	return r
}

// NewMember http parameter model
type NewMember struct {
	ServerID string `json:"server_id"`
	Address  string `json:"address"`
}

// JoinRequestParams exposed http model for join handler
type JoinRequestParams struct {
	NewMember *NewMember `json:"new_member"`
}

// Bind validate JoinRequestParams
func (j *JoinRequestParams) Bind(_ *http.Request) error {

	if j.NewMember == nil {
		return errors.New("invalid new_member parameter")
	}

	if len(j.NewMember.Address) < 5 {
		return errors.New("invalid address parameter")
	} else if len(j.NewMember.ServerID) < 16 {
		return errors.New("invalid server id parameter")
	}

	return nil
}

// Member http response model
type Member struct {
	Index    uint64 `json:"index"`
	ServerID string `json:"server_id"`
	Address  string `json:"address"`
}

// JoinResponse exposed http response model for join handler
type JoinResponse struct {
	Member  *Member `json:"member"`
	Elapsed int64   `json:"elapsed"`
}

// NewJoinResponse constructor
func NewJoinResponse(index uint64, serverID, address string) *JoinResponse {
	return &JoinResponse{
		Member: &Member{
			Index:    index,
			ServerID: serverID,
			Address:  address,
		},
	}
}

// Render join response
func (jr *JoinResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	jr.Elapsed = 10
	return nil
}

// Join raft handler
func (rc *Raft) Join(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var params JoinRequestParams
	err := render.Bind(r, &params)
	if err != nil {
		rc.log.Errorf("unable to bind request to model %v", err)

		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			rc.log.Errorf("failed to render invalid request %v", err)
			_ = render.Render(w, r, ErrRender(err))
		}

		return
	}

	m, err := rc.service.Join(ctx, &domain.NewMember{
		ServerID: params.NewMember.ServerID,
		Addr:     params.NewMember.Address,
	})
	if err != nil {
		rc.log.Errorf("new member unable to join raft %v", err)
		http.Error(w, "failed to join", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewJoinResponse(m.Index, m.ServerID, m.Addr))
	if err != nil {
		rc.log.Errorf("unable to render response %v", err)
		_ = render.Render(w, r, ErrRender(err))
	}
}
