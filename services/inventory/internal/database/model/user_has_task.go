package model

import (
	"time"
)

type UserHasTask struct {
	UserID    int
	TaskID    int
	CreatedAt time.Time
}
