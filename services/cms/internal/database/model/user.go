package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        int
	RoleID    int
	PublicID  uuid.UUID
	Email     string
	Username  string
	FirstName string
	LastName  string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
