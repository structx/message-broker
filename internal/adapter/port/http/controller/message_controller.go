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
	g.GET("/messages/:hash", m.fetchMessage)
	g.GET("/messages", m.listMessages)
	g.GET("/messages/topics", m.listTopics)
	g.GET("/messages/topics/:topic", m.listByTopic)
}

// MessagePayload http message model
type MessagePayload struct {
	Hash      string    `json:"hash"`
	Topic     string    `json:"topic"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// GetMessageResponse http get message response model
type GetMessageResponse struct {
	Payload *MessagePayload `json:"payload"`
}

func (m *Messages) fetchMessage(c echo.Context) error {

	h := c.Param("hash")

	msg, err := m.m.Read(h)
	if err != nil {
		m.log.Errorf("failed to read message %v", err)
		return c.String(http.StatusInternalServerError, "unable to read message")
	}

	response := &GetMessageResponse{
		Payload: &MessagePayload{
			Hash:      msg.Hash,
			Topic:     msg.Topic,
			Payload:   msg.Payload,
			CreatedAt: msg.CreatedAt,
		},
	}

	return c.JSON(http.StatusAccepted, response)
}

// ListMessagesResponse http list messages response model
type ListMessagesResponse struct {
	Payload []*MessagePayload `json:"payload"`
	Elapsed int64             `json:"elapsed"`
}

func (m *Messages) listMessages(c echo.Context) error {

	start := time.Now()

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

	messageSlice, err := m.m.List(int(l64), int(o64))
	if err != nil {
		m.log.Errorf("failed to list messages %v", err)
		return c.String(http.StatusInternalServerError, "failed to list messages")
	}

	payload := make([]*MessagePayload, 0, len(messageSlice))

	for _, msg := range messageSlice {
		payload = append(payload, &MessagePayload{
			Hash:      msg.Hash,
			Topic:     msg.Topic,
			Payload:   msg.Payload,
			CreatedAt: msg.CreatedAt,
		})
	}

	response := &ListMessagesResponse{
		Payload: payload,
		Elapsed: time.Now().Sub(start).Milliseconds(),
	}

	return c.JSON(http.StatusAccepted, response)
}

// ListTopicsResponse http list topics response model
type ListTopicsResponse struct {
	Topics []string `json:"topics"`
}

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

// PartialMessagePayload http partial message model
type PartialMessagePayload struct {
	Hash  string `json:"hash"`
	Topic string `json:"topic"`
}

// ListByTopicResponse http list by topic response model
type ListByTopicResponse struct {
	Elapsed int64                    `json:"elapsed"`
	Payload []*PartialMessagePayload `json:"payload"`
}

func (m *Messages) listByTopic(c echo.Context) error {

	start := time.Now()

	l := c.QueryParam("limit")
	o := c.QueryParam("offset")

	topic := c.Param("topic")
	if topic == "" {
		return c.String(http.StatusBadRequest, "invalid topic parameter")
	}

	l64, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid limit query parameter")
	}

	o64, err := strconv.ParseInt(o, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid offset query parameter")
	}

	partialMessageSlice, err := m.m.ListByTopic(topic, int(l64), int(o64))
	if err != nil {
		m.log.Errorf("failed to list messages by topic %v", err)
		return c.String(http.StatusInternalServerError, "failed to list messages by topic")
	}

	messageSlice := make([]*PartialMessagePayload, 0, len(partialMessageSlice))

	for _, p := range partialMessageSlice {
		messageSlice = append(messageSlice, &PartialMessagePayload{
			Hash:  p.Hash,
			Topic: p.Topic,
		})
	}

	response := &ListByTopicResponse{
		Elapsed: time.Now().Sub(start).Milliseconds(),
		Payload: messageSlice,
	}

	return c.JSON(http.StatusAccepted, response)
}
