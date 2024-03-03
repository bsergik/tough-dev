package model

import (
	"time"
)

type TaskHasStatus struct {
	TaskID    int
	StatusID  int
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
