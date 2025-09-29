package tasks

import "time"

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

const (
	TodoStr       = "To Do"
	InProgressStr = "In Progress"
	DoneStr       = "Done"
)

type Task struct {
	Id            int
	Description   string
	CurrentStatus Status
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TaskStorage interface {
	Save(tasks []Task) error
	Load() ([]Task, error)
}
