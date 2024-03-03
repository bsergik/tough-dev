package model

type TaskStatus int

const (
	TaskStatusUnknown    TaskStatus = 0
	TaskStatusCreated    TaskStatus = 1
	TaskStatusAssigned   TaskStatus = 2
	TaskStatusInProgress TaskStatus = 3
	TaskStatusClosed     TaskStatus = 4
	TaskStatusDone       TaskStatus = 5
)
