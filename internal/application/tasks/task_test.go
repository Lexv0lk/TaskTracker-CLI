//go:generate mockgen -destination=mocks/storage.go -package=mocks  github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks TaskStorage
package tasks

import (
	"fmt"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks/mocks"
	domain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		addingTask    domain.Task
		testStorageFn func(t *testing.T, addingTask domain.Task) domain.TaskStorage
		expectedErr   error
	}

	tests := []testCase{
		{
			name: "Successful Add",
			addingTask: domain.Task{
				Id:            1,
				Description:   "New Task",
				CurrentStatus: domain.Todo,
				CreatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, addingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)

				firstCall := result.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)
				result.EXPECT().Save(gomock.Eq([]domain.Task{addingTask})).Times(1).After(firstCall)

				return result
			},
			expectedErr: nil,
		},
		{
			name: "Storage Load Error",
			addingTask: domain.Task{
				Description: "New Task",
			},
			testStorageFn: func(t *testing.T, addingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)

				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name:       "Storage Save Error",
			addingTask: domain.Task{},
			testStorageFn: func(t *testing.T, addingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)

				firstCall := result.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)
				result.EXPECT().Save(gomock.Any()).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name: "Add Task with Existing Tasks",
			addingTask: domain.Task{
				Id:            3,
				Description:   "Another Task",
				CurrentStatus: domain.Todo,
				CreatedAt:     time.Date(2025, 2, 2, 15, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2025, 2, 2, 15, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, addingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.InProgress, CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(existingTasks, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(append(existingTasks, addingTask))).Times(1).After(firstCall)

				return result
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskStorage := tt.testStorageFn(t, tt.addingTask)
			err := addTask(taskStorage, tt.addingTask.Description, func() time.Time {
				return tt.addingTask.CreatedAt
			})

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		updatingTask  domain.Task
		testStorageFn func(t *testing.T, updatingTask domain.Task) domain.TaskStorage
		expectedErr   error
	}

	tests := []testCase{
		{
			name: "Successful Update",
			updatingTask: domain.Task{
				Id:            1,
				Description:   "Updated Task",
				CurrentStatus: domain.Todo,
				CreatedAt:     time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Old Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}
				expectedTasks := []domain.Task{updatingTask}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(existingTasks, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(expectedTasks)).Times(1).After(firstCall)

				return result
			},
			expectedErr: nil,
		},
		{
			name: "Task Not Found",
			updatingTask: domain.Task{
				Id: 2,
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Some Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return(existingTasks, nil).Times(1)

				return result
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 2),
		},
		{
			name:         "Storage Load Error",
			updatingTask: domain.Task{},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name:         "Storage Save Error",
			updatingTask: domain.Task{Id: 1, Description: "Updated Task", UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				tasks := []domain.Task{
					{Id: 1, Description: "Some Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(tasks, nil).Times(1)
				result.EXPECT().Save(gomock.Any()).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name: "Update Task Among Multiple Tasks",
			updatingTask: domain.Task{
				Id:            2,
				Description:   "Updated Task 2",
				CurrentStatus: domain.InProgress,
				CreatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.InProgress, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
					{Id: 3, Description: "Task 3", CurrentStatus: domain.Done, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC)},
				}
				expectedTasks := []domain.Task{
					existingTasks[0],
					updatingTask,
					existingTasks[2],
				}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(existingTasks, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(expectedTasks)).Times(1).After(firstCall)

				return result
			},
		},
		{
			name:         "Empty task list",
			updatingTask: domain.Task{Id: 5},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)

				return result
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 5),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskStorage := tt.testStorageFn(t, tt.updatingTask)
			err := updateTask(taskStorage, tt.updatingTask.Id, tt.updatingTask.Description, func() time.Time {
				return tt.updatingTask.UpdatedAt
			})

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		updatingTask  domain.Task
		newStatus     domain.Status
		testStorageFn func(t *testing.T, updatingTask domain.Task) domain.TaskStorage
		expectedErr   error
	}

	tests := []testCase{
		{
			name: "Successful Update",
			updatingTask: domain.Task{
				Id:            1,
				Description:   "Old Task",
				CurrentStatus: domain.InProgress,
				CreatedAt:     time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Old Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}
				expectedTasks := []domain.Task{updatingTask}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(existingTasks, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(expectedTasks)).Times(1).After(firstCall)

				return result
			},
			expectedErr: nil,
		},
		{
			name: "Task Not Found",
			updatingTask: domain.Task{
				Id: 2,
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Some Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return(existingTasks, nil).Times(1)

				return result
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 2),
		},
		{
			name:         "Storage Load Error",
			updatingTask: domain.Task{},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name:         "Storage Save Error",
			updatingTask: domain.Task{Id: 1, Description: "Some Task", CurrentStatus: domain.Done, UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				tasks := []domain.Task{
					{Id: 1, Description: "Some Task", CurrentStatus: domain.Todo, CreatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(tasks, nil).Times(1)
				result.EXPECT().Save(gomock.Any()).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expectedErr: assert.AnError,
		},
		{
			name: "Update Task Among Multiple Tasks",
			updatingTask: domain.Task{
				Id:            2,
				Description:   "Task 2",
				CurrentStatus: domain.Done,
				CreatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC),
			},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.InProgress, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
					{Id: 3, Description: "Task 3", CurrentStatus: domain.Done, CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC)},
				}
				expectedTasks := []domain.Task{
					existingTasks[0],
					updatingTask,
					existingTasks[2],
				}

				result := mocks.NewMockTaskStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(existingTasks, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(expectedTasks)).Times(1).After(firstCall)

				return result
			},
		},
		{
			name:         "Empty task list",
			updatingTask: domain.Task{Id: 5},
			testStorageFn: func(t *testing.T, updatingTask domain.Task) domain.TaskStorage {
				t.Helper()

				result := mocks.NewMockTaskStorage(ctrl)
				result.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)

				return result
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 5),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskStorage := tt.testStorageFn(t, tt.updatingTask)
			err := updateTaskStatus(taskStorage, tt.updatingTask.Id, tt.updatingTask.CurrentStatus, func() time.Time {
				return tt.updatingTask.UpdatedAt
			})

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		taskID        int
		testStorageFn func(t *testing.T, taskID int) domain.TaskStorage
		expectedErr   error
	}

	tests := []testCase{
		{
			name:   "Successful Delete",
			taskID: 1,
			testStorageFn: func(t *testing.T, taskID int) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1"},
					{Id: 2, Description: "Task 2"},
				}
				expectedTasks := []domain.Task{
					existingTasks[1], // Task 1 удалена
				}

				mock := mocks.NewMockTaskStorage(ctrl)
				firstCall := mock.EXPECT().Load().Return(existingTasks, nil).Times(1)
				mock.EXPECT().Save(gomock.Eq(expectedTasks)).Times(1).After(firstCall)

				return mock
			},
			expectedErr: nil,
		},
		{
			name:   "Task Not Found",
			taskID: 3,
			testStorageFn: func(t *testing.T, taskID int) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1"},
					{Id: 2, Description: "Task 2"},
				}

				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return(existingTasks, nil).Times(1)

				return mock
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 3),
		},
		{
			name:   "Storage Load Error",
			taskID: 1,
			testStorageFn: func(t *testing.T, taskID int) domain.TaskStorage {
				t.Helper()

				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return mock
			},
			expectedErr: assert.AnError,
		},
		{
			name:   "Storage Save Error",
			taskID: 1,
			testStorageFn: func(t *testing.T, taskID int) domain.TaskStorage {
				t.Helper()

				existingTasks := []domain.Task{
					{Id: 1, Description: "Task 1"},
					{Id: 2, Description: "Task 2"},
				}

				mock := mocks.NewMockTaskStorage(ctrl)
				firstCall := mock.EXPECT().Load().Return(existingTasks, nil).Times(1)
				mock.EXPECT().Save(gomock.Any()).Return(assert.AnError).Times(1).After(firstCall)

				return mock
			},
			expectedErr: assert.AnError,
		},
		{
			name:   "Empty task list",
			taskID: 1,
			testStorageFn: func(t *testing.T, taskID int) domain.TaskStorage {
				t.Helper()

				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)

				return mock
			},
			expectedErr: fmt.Errorf("task with id [%d] not found", 1),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskStorage := tt.testStorageFn(t, tt.taskID)
			err := deleteTask(taskStorage, tt.taskID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllTasksList(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		mockTasks     []domain.Task
		testStorageFn func(t *testing.T) domain.TaskStorage
		expectedOut   string
		expectedErr   error
	}

	tests := []testCase{
		{
			name: "Multiple tasks",
			mockTasks: []domain.Task{
				{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
				{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
			},
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
				}, nil).Times(1)
				return mock
			},
			expectedOut: fmt.Sprintf(
				"%-3s %-20s %-12s %-16s %-16s\n%s%s",
				"ID", "Description", "Status", "Created At", "Updated At",
				getTaskShortDescription(domain.Task{
					Id: 1, Description: "Task 1", CurrentStatus: domain.Todo,
					CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC),
				}),
				getTaskShortDescription(domain.Task{
					Id: 2, Description: "Task 2", CurrentStatus: domain.Done,
					CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC),
				}),
			),
			expectedErr: nil,
		},
		{
			name:      "Empty task list",
			mockTasks: []domain.Task{},
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)
				return mock
			},
			expectedOut: fmt.Sprintf("%-3s %-20s %-12s %-16s %-16s\n",
				"ID", "Description", "Status", "Created At", "Updated At"),
			expectedErr: nil,
		},
		{
			name: "Storage Load Error",
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return(nil, assert.AnError).Times(1)
				return mock
			},
			expectedOut: "",
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskStorage := tt.testStorageFn(t)
			out, err := getAllTasksList(taskStorage)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, out)
			}
		})
	}
}

func TestGetFilteredTasksList(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name          string
		filterStatus  domain.Status
		mockTasks     []domain.Task
		testStorageFn func(t *testing.T) domain.TaskStorage
		expectedOut   string
		expectedErr   error
	}

	tests := []testCase{
		{
			name:         "Filter Todo Tasks",
			filterStatus: domain.Todo,
			mockTasks: []domain.Task{
				{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
				{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
			},
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
				}, nil).Times(1)
				return mock
			},
			expectedOut: fmt.Sprintf(
				"%s%s",
				getTaskListHeader(),
				getTaskShortDescription(domain.Task{
					Id: 1, Description: "Task 1", CurrentStatus: domain.Todo,
					CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC),
				}),
			),
			expectedErr: nil,
		},
		{
			name:         "No tasks match filter",
			filterStatus: domain.InProgress,
			mockTasks: []domain.Task{
				{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
				{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
			},
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{
					{Id: 1, Description: "Task 1", CurrentStatus: domain.Todo, CreatedAt: time.Date(2025, 9, 29, 12, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 29, 12, 30, 0, 0, time.UTC)},
					{Id: 2, Description: "Task 2", CurrentStatus: domain.Done, CreatedAt: time.Date(2025, 9, 28, 9, 15, 0, 0, time.UTC), UpdatedAt: time.Date(2025, 9, 28, 17, 45, 0, 0, time.UTC)},
				}, nil).Times(1)
				return mock
			},
			expectedOut: getTaskListHeader(), // только заголовок
			expectedErr: nil,
		},
		{
			name:         "Storage Load Error",
			filterStatus: domain.Todo,
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return(nil, assert.AnError).Times(1)
				return mock
			},
			expectedOut: "",
			expectedErr: assert.AnError,
		},
		{
			name:         "Empty task list",
			filterStatus: domain.Todo,
			testStorageFn: func(t *testing.T) domain.TaskStorage {
				t.Helper()
				mock := mocks.NewMockTaskStorage(ctrl)
				mock.EXPECT().Load().Return([]domain.Task{}, nil).Times(1)
				return mock
			},
			expectedOut: getTaskListHeader(),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := tt.testStorageFn(t)
			out, err := getFilteredTasksList(storage, tt.filterStatus)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, out)
			}
		})
	}
}
