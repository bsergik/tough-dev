package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Task struct {
	ID          int
	PublicID    uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
