package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"
)

// Tx blockchain level transaction model
type Tx struct {
	ID        []byte `json:"id"`
	Topic     string `json:"topic"`
	Action    string `json:"pattern"`
	Payload   []byte `json:"payload"`
	Timestamp string `json:"timestamp"`
	Sig       string `json:"sig"`
}

// NewTx return new transaction
func NewTx(data, action string, payload []byte, signature string) *Tx {
	return &Tx{
		Topic:     data,
		Action:    action,
		Payload:   payload,
		Timestamp: time.Now().String(),
		Sig:       signature,
	}
}

// SetID set transaction hash
func (t *Tx) SetID() {
	h := sha3.New224()
	h.Write([]byte(fmt.Sprintf("%s:%s:%x", t.Topic, t.Action, t.Sig)))
	t.ID = h.Sum(nil)
}
