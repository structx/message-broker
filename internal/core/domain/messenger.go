package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NewMessage
type NewMessage struct {
	Topic     string
	Payload   []byte
	Publisher string
}

// Message
type Message struct {
	ID        string
	Topic     string
	Payload   []byte
	CreatedAt time.Time
}

// Messenger
//
//go:generate mockery --name Messenger
type Messenger interface {
	// Create
	Create(ctx context.Context, newMessage *NewMessage) (*Message, error)
	Read(ctx context.Context, UID uuid.UUID) (*Message, error)
	List(ctx context.Context, limit, offset int) ([]*Message, error)
	ListByTopic(ctx context.Context, topic string, limit, offset int) ([]*Message, error)
	ListTopics(ctx context.Context) ([]string, error)
}
