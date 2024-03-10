package model

import (
	"time"
)

type FailedMessage struct {
	ID         int
	Topic      string
	GroupID    string
	Message    []byte
	ReceivedAt time.Time
	CreatedAt  time.Time
}
