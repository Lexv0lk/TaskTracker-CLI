//go:generate mockgen -destination=mocks/files.go -package=mocks io WriteCloser,ReadCloser
package task

import (
	"encoding/json"
	"fmt"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/task/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type TestCase struct {
		name          string
		writeCloserFn func(t *testing.T, tasks []task) io.WriteCloser
		tasks         []task
		expectedErr   error
	}

	tests := []TestCase{
		{
			name: "Successful Save",
			writeCloserFn: func(t *testing.T, tasks []task) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(tasks, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			tasks: []task{
				{Id: 1, Description: "Task 1", CurrentStatus: todo},
				{Id: 2, Description: "Task 2", CurrentStatus: inProgress},
				{Id: 3, Description: "Task 3", CurrentStatus: done},
			},
			expectedErr: nil,
		},
		{
			name: "Error writer",
			writeCloserFn: func(t *testing.T, tasks []task) io.WriteCloser {
				t.Helper()

				testErr := fmt.Errorf("test error")

				result := mocks.NewMockWriteCloser(ctrl)
				result.EXPECT().Write(gomock.Any()).Return(0, testErr).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			tasks:       []task{},
			expectedErr: fmt.Errorf("test error"),
		},
		{
			name: "Empty Task List",
			writeCloserFn: func(t *testing.T, tasks []task) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(tasks, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			tasks:       []task{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := saveToFile(tt.writeCloserFn(t, tt.tasks), tt.tasks)

			if tt.expectedErr != nil {
				assert.EqualError(err, tt.expectedErr.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}
