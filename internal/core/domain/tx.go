package domain

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Tx
type Tx struct {
	ID        []byte `json:"id"`
	Topic     string `json:"topic"`
	Pattern   string `json:"pattern"`
	Payload   []byte `json:"payload"`
	Timestamp string `json:"timestamp"`
	Sig       []byte `json:"sig"`
}

// SetID
func (t *Tx) SetID() {
	h := sha3.New224()
	h.Write([]byte(fmt.Sprintf("%s:%s:%x", t.Topic, t.Pattern, t.Sig)))
	t.ID = h.Sum(nil)
}
