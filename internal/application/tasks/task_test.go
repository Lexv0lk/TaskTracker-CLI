//go:generate mockgen -destination=mocks/storage.go -package=mocks  github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks TaskStorage
package tasks

import (
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
