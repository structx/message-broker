package application

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/trevatk/block-broker/internal/core/domain"
)

// MessagingService
type MessagingService struct {
	c domain.Chain
}

// interface compliance
var _ domain.Messenger = (*MessagingService)(nil)

// NewMessagingService
func NewMessagingService(c domain.Chain) *MessagingService {
	return &MessagingService{
		c: c,
	}
}

// Create
func (m *MessagingService) Create(ctx context.Context, newMessage *domain.NewMessage) (*domain.Message, error) {

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

// Read
func (m *MessagingService) Read(ctx context.Context, UID uuid.UUID) (*domain.Message, error) {
	return nil, nil
}

// Acknowledge
func (m *MessagingService) Acknowledge(ctx context.Context) error {
	return nil
}

// List
func (m *MessagingService) List(ctx context.Context, limit, offset int) ([]*domain.Message, error) {

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

// ListByTopic
func (m *MessagingService) ListByTopic(ctx context.Context, topic string, limit, offset int) ([]*domain.Message, error) {

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

// ListTopics
func (m *MessagingService) ListTopics(ctx context.Context) ([]string, error) {

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
