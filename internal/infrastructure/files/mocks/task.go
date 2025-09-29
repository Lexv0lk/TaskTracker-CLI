package mocks

import "time"

type StatusMock int

const (
	Todo StatusMock = iota
	InProgress
	Done
)

type TaskMock struct {
	Id            int
	Description   string
	CurrentStatus StatusMock
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
