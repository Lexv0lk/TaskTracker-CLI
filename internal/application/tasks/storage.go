package tasks

import (
	domain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/infrastructure/files"
)

type taskFileStorage struct {
}

var defaultTaskStorage domain.TaskStorage = &taskFileStorage{}

func (t *taskFileStorage) Save(tasks []domain.Task) error {
	return files.SaveToFile(tasks)
}

func (t *taskFileStorage) Load() ([]domain.Task, error) {
	return files.GetFromFile[[]domain.Task]()
}
