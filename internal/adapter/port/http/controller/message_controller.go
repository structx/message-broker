package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/trevatk/block-broker/internal/core/domain"
)

const (
	messageId = "/message/{id}"
)

// MessageController
type MessageController struct {
	log *zap.SugaredLogger
	m   domain.Messenger
}

// interface compliance
var _ Controller = (*MessageController)(nil)

// NewMessageController
func NewMessageController(logger *zap.Logger, messenger domain.Messenger) *MessageController {
	return &MessageController{
		log: logger.Sugar().Named("message_controller"),
		m:   messenger,
	}
}

// RegisterRoutesV1
func (mc *MessageController) RegisterRoutesV1(s *fuego.Server) {
	fuego.GetStd(s, "/api/v1/message/{id}", mc.Get).
		Summary("Fetch message").
		Description("Fetch message by id").
		AddTags("Messages").
		OperationID("fetchMessage").
		QueryParam("messageId", "message unique identifier", fuego.OpenAPIParam{
			Required: true,
			Type:     "string",
			Example:  "322befbd-ed13-4566-93e7-24fe87e5306f",
		})
	fuego.GetStd(s, "/api/v1/messages/topics", mc.ListTopics).
		Summary("List topics").
		Description("List all topics").
		AddTags("Topics").
		OperationID("listTopics")
}

// MessagePayload
type MessagePayload struct {
	UID       string    `json:"uid"`
	Topic     string    `json:"topic"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// GetMessageResponse
type GetMessageResponse struct {
	Payload *MessagePayload `json:"payload"`
}

// Get
func (mc *MessageController) Get(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.String()
	uris := strings.Split(uri, "/")

	id := uris[len(uris)-1]

	uid, err := uuid.Parse(id)
	if err != nil {
		mc.log.Errorf("failed to parse uuid %v", err)
		http.Error(w, "invalid request parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	msg, err := mc.m.Read(ctx, uid)
	if err != nil {
		mc.log.Errorf("failed to read message %v", err)
		http.Error(w, "unable to read message", http.StatusInternalServerError)
		return
	}

	response := &GetMessageResponse{
		Payload: &MessagePayload{
			UID:       msg.ID,
			Topic:     msg.Topic,
			Payload:   msg.Payload,
			CreatedAt: msg.CreatedAt,
		},
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		mc.log.Errorf("failed to encode response %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

// ListTopicsResponse
type ListTopicsReponse struct {
	Topics []string `json:"topics"`
}

// ListTopics
func (mc *MessageController) ListTopics(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	topics, err := mc.m.ListTopics(ctx)
	if err != nil {
		mc.log.Errorf("m.ListTopics: %v", err)
		http.Error(w, "unable to list topics", http.StatusInternalServerError)
		return

	}

	response := &ListTopicsReponse{
		Topics: topics,
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		mc.log.Errorf("json.NewEncoder: %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

// ListByTopic
func (mc *MessageController) ListByTopic(w http.ResponseWriter, r *http.Request) {}
