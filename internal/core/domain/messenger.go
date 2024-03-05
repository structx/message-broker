package domain

import (
	"time"
)

// PayloadAction domain actions
type PayloadAction string

const (
	// Publish message
	Publish PayloadAction = "publish"
	// Subscribe to topic
	Subscribe PayloadAction = "subscribe"
)

// NewMessage ...
type NewMessage struct {
	Topic     string
	Payload   []byte
	Publisher string
}

// Message ...
type Message struct {
	ID        string
	Topic     string
	Payload   []byte
	CreatedAt time.Time
}

// Messenger messaging service interface
//
//go:generate mockery --name Messenger
type Messenger interface {
	// Create message
	Create(newMessage *NewMessage) (*Message, error)
	// Read message by hash
	Read(hash string) (*Message, error)
	// List messages
	List(limit, offset int) ([]*Message, error)
	// ListByTopic messages by topic
	ListByTopic(topic string, limit, offset int) ([]*Message, error)
	// ListTopics retrieve message topics
	ListTopics(limit, offset int) ([]string, error)
}
