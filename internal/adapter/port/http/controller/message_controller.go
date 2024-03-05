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

// Messages controller class
type Messages struct {
	log *zap.SugaredLogger
	m   domain.Messenger
}

// interface compliance
var _ Controller = (*Messages)(nil)

// NewMessages return new message controller
func NewMessages(logger *zap.Logger, messenger domain.Messenger) *Messages {
	return &Messages{
		log: logger.Sugar().Named("message_controller"),
		m:   messenger,
	}
}

// RegisterRoutesV1 register routes on v1 router
func (m *Messages) RegisterRoutesV1(s *fuego.Server) {
	fuego.GetStd(s, "/api/v1/message/{hash}", m.Get).
		Summary("Fetch message").
		Description("Fetch message by hash").
		AddTags("Messages").
		OperationID("fetchMessage").
		QueryParam("hash", "message hash", fuego.OpenAPIParam{
			Required: true,
			Type:     "string",
			Example:  "0006fe63d8b226c08bb7ce6dc7e0f2beb6436bcb8531184c0656dbb1",
		})
	fuego.GetStd(s, "/api/v1/messages/topics", m.ListTopics).
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
func (m *Messages) Get(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	urlslice := strings.Split(url, "/")

	hash := urlslice[len(urlslice)-1]

	msg, err := m.m.Read(hash)
	if err != nil {
		m.log.Errorf("failed to read message %v", err)
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
		m.log.Errorf("failed to encode response %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

// ListTopicsResponse http list topics response model
type ListTopicsResponse struct {
	Topics []string `json:"topics"`
}

// ListTopics fetch topics
func (m *Messages) ListTopics(w http.ResponseWriter, _ *http.Request) {

	topics, err := m.m.ListTopics(0, 0)
	if err != nil {
		m.log.Errorf("m.ListTopics: %v", err)
		http.Error(w, "unable to list topics", http.StatusInternalServerError)
		return

	}

	response := &ListTopicsResponse{
		Topics: topics,
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		m.log.Errorf("json.NewEncoder: %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}

// ListByTopic fetch messages by topic
func (m *Messages) ListByTopic(_ http.ResponseWriter, _ *http.Request) {
	// TODO:
	// implement handler
}
