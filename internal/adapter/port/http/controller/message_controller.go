package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

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
func (m *Messages) RegisterRoutesV1(g *echo.Group) {
	g.GET("/message/:hash", m.fetchMessage)
	g.GET("/message/topics", m.listTopics)
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
func (m *Messages) fetchMessage(c echo.Context) error {

	h := c.Param("hash")

	msg, err := m.m.Read(h)
	if err != nil {
		m.log.Errorf("failed to read message %v", err)
		return c.String(http.StatusInternalServerError, "unable to read message")
	}

	response := &GetMessageResponse{
		Payload: &MessagePayload{
			UID:       msg.ID,
			Topic:     msg.Topic,
			Payload:   msg.Payload,
			CreatedAt: msg.CreatedAt,
		},
	}

	return c.JSON(http.StatusAccepted, response)
}

// ListTopicsResponse http list topics response model
type ListTopicsResponse struct {
	Topics []string `json:"topics"`
}

// ListTopics fetch topics
func (m *Messages) listTopics(c echo.Context) error {

	l := c.QueryParam("limit")
	o := c.QueryParam("offset")

	l64, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid limit query parameter")
	}

	o64, err := strconv.ParseInt(o, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid offset query parameter")
	}

	topics, err := m.m.ListTopics(int(l64), int(o64))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to list topics")
	}

	response := &ListTopicsResponse{
		Topics: topics,
	}

	return c.JSON(http.StatusAccepted, response)
}

// ListByTopic fetch messages by topic
func (m *Messages) ListByTopic(_ http.ResponseWriter, _ *http.Request) {
	// TODO:
	// implement handler
}
