// Package application service logic
package application

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/trevatk/block-broker/internal/core/domain"
)

// MessagingService messaging service implementation
type MessagingService struct {
	c domain.Chain
}

// interface compliance
var _ domain.Messenger = (*MessagingService)(nil)

// NewMessagingService return new messaging service class
func NewMessagingService(c domain.Chain) *MessagingService {
	return &MessagingService{
		c: c,
	}
}

// Create transaction
func (m *MessagingService) Create(newMessage *domain.NewMessage) (*domain.Message, error) {

	action := string(domain.Publish)
	hash, err := m.c.AddTx(newMessage.Topic, action, newMessage.Payload, newMessage.Publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to add tranaction %v", err)
	}

	return &domain.Message{
		ID:        hash,
		Topic:     newMessage.Topic,
		Payload:   newMessage.Payload,
		CreatedAt: time.Now(),
	}, nil
}

// Read transaction by hash
func (m *MessagingService) Read(_ string) (*domain.Message, error) {
	// TODO:
	// implement logic
	return nil, nil
}

// Acknowledge todo
func (m *MessagingService) Acknowledge(_ context.Context) error {
	// TODO:
	// implement logic
	return nil
}

// List messages by limit and offset
func (m *MessagingService) List(limit, offset int) ([]*domain.Message, error) {

	txslice, err := m.c.ListTransactions(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by action %v", err)
	}

	messageslice := make([]*domain.Message, 0, limit)

	for _, tx := range txslice {

		timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tx timestamp %v", err)
		}

		messageslice = append(messageslice, &domain.Message{
			ID:        hex.EncodeToString(tx.ID),
			Topic:     tx.Topic,
			Payload:   tx.Payload,
			CreatedAt: timestamp,
		})
	}

	return messageslice, nil
}

// ListByTopic retrieve messages by topic
func (m *MessagingService) ListByTopic(topic string, limit, offset int) ([]*domain.Message, error) {

	txslice, err := m.c.ListTransactionsByAction(topic, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by action %v", err)
	}

	messageslice := make([]*domain.Message, 0, limit)

	for _, tx := range txslice {

		timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tx timestamp %v", err)
		}

		messageslice = append(messageslice, &domain.Message{
			ID:        hex.EncodeToString(tx.ID),
			Topic:     tx.Topic,
			Payload:   tx.Payload,
			CreatedAt: timestamp,
		})
	}

	return messageslice, nil
}

// ListTopics retrieve all topics from chain
func (m *MessagingService) ListTopics(limit, offset int) ([]string, error) {

	txslice, err := m.c.ListTransactions(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions %v", err)
	}

	topicSlice := make([]string, 0, limit)

	for _, tx := range txslice {
		topicSlice = append(topicSlice, tx.Topic)
	}

	return topicSlice, nil
}
