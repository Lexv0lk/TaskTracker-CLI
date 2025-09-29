package tasks

import (
	"fmt"
	domain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"sort"
	"strings"
	"time"
)

func AddTask(description string) (domain.Task, error) {
	return addTask(defaultTaskStorage, description, time.Now)
}

func UpdateTask(id int, description string) error {
	return updateTask(defaultTaskStorage, id, description, time.Now)
}

func UpdateTaskStatus(id int, status domain.Status) error {
	return updateTaskStatus(defaultTaskStorage, id, status, time.Now)
}

func DeleteTask(id int) error {
	return deleteTask(defaultTaskStorage, id)
}

func GetAllTasks() (string, error) {
	return getAllTasksList(defaultTaskStorage)
}

func GetTasks(status domain.Status) (string, error) {
	return getFilteredTasksList(defaultTaskStorage, status)
}

func ParseStatusString(statusStr string) (domain.Status, error) {
	switch strings.ToLower(statusStr) {
	case domain.TodoStr:
		return domain.Todo, nil
	case domain.InProgressStr:
		return domain.InProgress, nil
	case domain.DoneStr:
		return domain.Done, nil
	default:
		return 0, fmt.Errorf("invalid status string: %s", statusStr)
	}
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
