package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-fuego/fuego"

	"go.uber.org/zap"

	"github.com/trevatk/block-broker/internal/core/domain"
)

// MessageController message controller class
type MessageController struct {
	log *zap.SugaredLogger
	m   domain.Messenger
}

// interface compliance
var _ Controller = (*MessageController)(nil)

// NewMessageController return new message controller
func NewMessageController(logger *zap.Logger, messenger domain.Messenger) *MessageController {
	return &MessageController{
		log: logger.Sugar().Named("message_controller"),
		m:   messenger,
	}
}

// RegisterRoutesV1 register routes on v1 router
func (mc *MessageController) RegisterRoutesV1(s *fuego.Server) {
	fuego.GetStd(s, "/api/v1/message/{id}", mc.Get).
		Summary("Fetch message").
		Description("Fetch message by hash").
		AddTags("Messages").
		OperationID("fetchMessage").
		QueryParam("messageID", "message hash", fuego.OpenAPIParam{
			Required: true,
			Type:     "string",
			Example:  "",
		})
	fuego.GetStd(s, "/api/v1/messages/topics", mc.ListTopics).
		Summary("List topics").
		Description("List all topics").
		AddTags("Topics").
		OperationID("listTopics")
}

// MessagePayload http message model
type MessagePayload struct {
	UID       string    `json:"uid"`
	Topic     string    `json:"topic"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// GetMessageResponse http get message response model
type GetMessageResponse struct {
	Payload *MessagePayload `json:"payload"`
}

// Get fetch message by hash
func (mc *MessageController) Get(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	urlslice := strings.Split(url, "/")

	hash := urlslice[len(urlslice)-1]

	msg, err := mc.m.Read(hash)
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

// ListTopicsResponse http list topics response model
type ListTopicsResponse struct {
	Topics []string `json:"topics"`
}

// ListTopics fetch topics
func (mc *MessageController) ListTopics(w http.ResponseWriter, _ *http.Request) {

	topics, err := mc.m.ListTopics(0, 0)
	if err != nil {
		mc.log.Errorf("m.ListTopics: %v", err)
		http.Error(w, "unable to list topics", http.StatusInternalServerError)
		return

	}

	response := &ListTopicsResponse{
		Topics: topics,
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		mc.log.Errorf("json.NewEncoder: %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

// ListByTopic fetch messages by topic
func (mc *MessageController) ListByTopic(_ http.ResponseWriter, _ *http.Request) {
	// TODO:
	// implement handler
}
