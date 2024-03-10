package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Task struct {
	Message

	PublicID    uuid.UUID `json:"public_id"`
	UserID      uuid.UUID `json:"user_id"`
	StatusID    int       `json:"status_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
