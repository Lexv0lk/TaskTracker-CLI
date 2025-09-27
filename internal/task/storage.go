package task

import "github.com/Lexv0lk/TaskTracker-CLI/internal/files"

type TaskStorage interface {
	Save(tasks []task) error
	Load() ([]task, error)
}

type taskFileStorage struct {
}

var defaultTaskStorage TaskStorage = &taskFileStorage{}

func (t *taskFileStorage) Save(tasks []task) error {
	return files.SaveToFile(tasks)
}

func (t *taskFileStorage) Load() ([]task, error) {
	return files.GetFromFile[[]task]()
}
