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

	tx := &domain.Tx{
		Topic:   newMessage.Topic,
		Payload: newMessage.Payload,
		Sig:     []byte(newMessage.Publisher),
	}
	tx.SetID()

	tx, err := m.c.InsertTx(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction %v", err)
	}

	return &domain.Message{
		ID:      hex.EncodeToString(tx.ID),
		Topic:   tx.Topic,
		Payload: tx.Payload,
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
func (m *MessagingService) List(_, _ int) ([]*domain.Message, error) {

	// TODO:
	// implement using limit and offset
	// with blockchain iterator

	block, err := m.c.GetBlock(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get first block %v", err)
	}

	messageSlice := make([]*domain.Message, 0)

	for _, tx := range block.Txs {

		createdAt, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse transaction timestamp %v", err)
		}

		messageSlice = append(messageSlice, &domain.Message{
			ID:        string(tx.Payload),
			Topic:     tx.Topic,
			Payload:   tx.Payload,
			CreatedAt: createdAt,
		})
	}

	return messageSlice, nil
}

// ListByTopic retrieve messages by topic
func (m *MessagingService) ListByTopic(topic string, _, _ int) ([]*domain.Message, error) {

	// TODO:
	// implement using limit and offset
	// with blockchain iterator

	block, err := m.c.GetBlock(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get first block %v", err)
	}

	messageSlice := make([]*domain.Message, 0)

	for _, tx := range block.Txs {

		if tx.Topic != topic {
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse transaction timestamp %v", err)
		}

		messageSlice = append(messageSlice, &domain.Message{
			ID:        string(tx.Payload),
			Topic:     tx.Topic,
			Payload:   tx.Payload,
			CreatedAt: createdAt,
		})
	}

	return messageSlice, nil
}

// ListTopics retrieve all topics from chain
func (m *MessagingService) ListTopics(_, _ int) ([]string, error) {

	// TODO:
	// implement using limit and offset
	// with blockchain iterator

	block, err := m.c.GetBlock(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get first block %v", err)
	}

	topicSlice := make([]string, 0)

	for _, tx := range block.Txs {
		topicSlice = append(topicSlice, tx.Topic)
	}

	return topicSlice, nil
}
