package tasks

import (
	domain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"sort"
	"time"
)

func AddTask(description string) error {
	return addTask(defaultTaskStorage, description, time.Now)
}

func addTask(taskStorage domain.TaskStorage, description string, now func() time.Time) error {
	tasks, err := taskStorage.Load()

	if err != nil {
		return err
	}

	newTask := domain.Task{
		Id:            getNextId(tasks),
		Description:   description,
		CurrentStatus: domain.Todo,
		CreatedAt:     now(),
		UpdatedAt:     now(),
	}

	tasks = append(tasks, newTask)
	return taskStorage.Save(tasks)
}

func getNextId(tasks []domain.Task) int {
	if len(tasks) == 0 {
		return 1
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})

	return tasks[len(tasks)-1].Id + 1
}
