package tasks

import "time"

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

const (
	TodoStr       = "todo"
	InProgressStr = "in-progress"
	DoneStr       = "done"
	UnknownStr    = "unknown"
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

func (s Status) String() string {
	switch s {
	case Todo:
		return TodoStr
	case InProgress:
		return InProgressStr
	case Done:
		return DoneStr
	default:
		return UnknownStr
	}
}
