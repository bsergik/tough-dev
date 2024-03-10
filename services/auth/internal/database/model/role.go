package model

import (
	"time"
)

type RoleType int

const (
	RoleTypeAdmin RoleType = 1 + iota
	RoleTypeUser
)

type Role struct {
	ID        int
	Title     string
	CreatedAt time.Time
}
