package task

import (
	"sort"
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

func AddTask(description string) error {
	return addTask(defaultTaskStorage, description, time.Now)
}

func addTask(taskStorage TaskStorage, description string, now func() time.Time) error {
	tasks, err := taskStorage.Load[[]task]()

	if err != nil {
		return err
	}

	newTask := task{
		Id:            getNextId(tasks),
		Description:   description,
		CurrentStatus: todo,
		CreatedAt:     now(),
		UpdatedAt:     now(),
	}

	tasks = append(tasks, newTask)
	return taskStorage.Save(tasks)
}

func getNextId(tasks []task) int {
	if len(tasks) == 0 {
		return 1
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})

	return tasks[len(tasks)-1].Id + 1
}
