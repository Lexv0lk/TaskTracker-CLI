package task

import (
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

type task struct {
	Id            int
	Description   string
	CurrentStatus status
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
