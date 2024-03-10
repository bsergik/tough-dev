package model

import (
	"time"
)

type MQOffset struct {
	Offset    int64
	UpdatedAt time.Time
}
