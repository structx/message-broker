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
	Signature string
}

// Message ...
type Message struct {
	Hash      string
	Topic     string
	Payload   []byte
	CreatedAt time.Time
}

// PartialMessage ...
type PartialMessage struct {
	Hash  string
	Topic string
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
	List(limit, offset int) ([]*PartialMessage, error)
	// ListByTopic messages by topic
	ListByTopic(topic string, limit, offset int) ([]*PartialMessage, error)
	// ListTopics retrieve message topics
	ListTopics(limit, offset int) ([]string, error)
}
