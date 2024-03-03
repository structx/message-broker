package domain

import "time"

// NewSubscription
type NewSubscription struct {
	Topic, Ip, Hostname string
}

// Subscription
type Subscription struct {
	ID        int64
	Topic     string
	Ip        string
	Hostname  string
	CreatedAt time.Time
}
