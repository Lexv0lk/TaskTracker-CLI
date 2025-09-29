package tasks

import (
	"fmt"
	domain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"strings"
	"time"
)

func addTask(taskStorage domain.TaskStorage, description string, now func() time.Time) (domain.Task, error) {
	tasks, err := taskStorage.Load()

	if err != nil {
		return domain.Task{}, err
	}

	newTask := domain.Task{
		Id:            getNextId(tasks),
		Description:   description,
		CurrentStatus: domain.Todo,
		CreatedAt:     now(),
		UpdatedAt:     now(),
	}

	tasks = append(tasks, newTask)
	return newTask, taskStorage.Save(tasks)
}

func updateTask(taskStorage domain.TaskStorage, id int, description string, now func() time.Time) error {
	tasks, err := taskStorage.Load()

	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = now()
			return taskStorage.Save(tasks)
		}
	}

	return fmt.Errorf("task with id [%d] not found", id)
}

func updateTaskStatus(taskStorage domain.TaskStorage, id int, status domain.Status, now func() time.Time) error {
	tasks, err := taskStorage.Load()

	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].CurrentStatus = status
			tasks[i].UpdatedAt = now()
			return taskStorage.Save(tasks)
		}
	}

	return fmt.Errorf("task with id [%d] not found", id)
}

func deleteTask(taskStorage domain.TaskStorage, id int) error {
	tasks, err := taskStorage.Load()

	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return taskStorage.Save(tasks)
		}
	}

	return fmt.Errorf("task with id [%d] not found", id)
}

func getAllTasksList(storage domain.TaskStorage) (string, error) {
	tasks, err := storage.Load()

	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(getTaskListHeader())

	for _, task := range tasks {
		builder.WriteString(getTaskShortDescription(task))
	}

	return builder.String(), nil
}

func getFilteredTasksList(storage domain.TaskStorage, status domain.Status) (string, error) {
	tasks, err := storage.Load()

	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(getTaskListHeader())

	for _, task := range tasks {
		if task.CurrentStatus == status {
			builder.WriteString(getTaskShortDescription(task))
		}
	}

	return builder.String(), nil
}

func getTaskListHeader() string {
	return fmt.Sprintf("%-3s %-20s %-12s %-16s %-16s\n",
		"ID", "Description", "Status", "Created At", "Updated At")
}

func getTaskShortDescription(task domain.Task) string {
	return fmt.Sprintf("%-3d %-20s %-12s %-16s %-16s\n",
		task.Id,
		task.Description,
		task.CurrentStatus.String(),
		task.CreatedAt.Format("2006-01-02 15:04"),
		task.UpdatedAt.Format("2006-01-02 15:04"),
	)
}
