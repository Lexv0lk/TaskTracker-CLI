//go:generate mockgen -destination=mocks/files.go -package=mocks io WriteCloser,ReadCloser
package files

import (
	"encoding/json"
	"fmt"
	mocks2 "github.com/Lexv0lk/TaskTracker-CLI/internal/infrastructure/files/mocks"
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
		writeCloserFn func(t *testing.T, tasks []mocks2.TaskMock) io.WriteCloser
		tasks         []mocks2.TaskMock
		expectedErr   error
	}

	tests := []TestCase{
		{
			name: "Successful Save",
			writeCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(tasks, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks2.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			tasks: []mocks2.TaskMock{
				{Id: 1, Description: "Task 1", CurrentStatus: mocks2.Todo},
				{Id: 2, Description: "Task 2", CurrentStatus: mocks2.InProgress},
				{Id: 3, Description: "Task 3", CurrentStatus: mocks2.Done},
			},
			expectedErr: nil,
		},
		{
			name: "Error writer",
			writeCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.WriteCloser {
				t.Helper()

				testErr := fmt.Errorf("test error")

				result := mocks2.NewMockWriteCloser(ctrl)
				result.EXPECT().Write(gomock.Any()).Return(0, testErr).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			tasks:       []mocks2.TaskMock{},
			expectedErr: fmt.Errorf("test error"),
		},
		{
			name: "Empty Task List",
			writeCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(tasks, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks2.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			tasks:       []mocks2.TaskMock{},
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

func TestGetFromFile(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type TestCase struct {
		name          string
		readCloserFn  func(t *testing.T, tasks []mocks2.TaskMock) io.ReadCloser
		expectedTasks []mocks2.TaskMock
		expectedErr   error
	}

	tests := []TestCase{
		{
			name: "Successful Read",
			readCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.ReadCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(tasks, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks2.NewMockReadCloser(ctrl)
				firstCall := result.EXPECT().Read(gomock.Any()).DoAndReturn(
					func(p []byte) (n int, err error) {
						copy(p, correctJson)
						return len(correctJson), nil
					}).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			expectedTasks: []mocks2.TaskMock{
				{Id: 1, Description: "Task 1", CurrentStatus: mocks2.Todo},
				{Id: 2, Description: "Task 2", CurrentStatus: mocks2.InProgress},
				{Id: 3, Description: "Task 3", CurrentStatus: mocks2.Done},
			},
			expectedErr: nil,
		},
		{
			name: "Error reader",
			readCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.ReadCloser {
				t.Helper()

				testErr := fmt.Errorf("test error")

				result := mocks2.NewMockReadCloser(ctrl)
				result.EXPECT().Read(gomock.Any()).Return(0, testErr).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			expectedTasks: nil,
			expectedErr:   fmt.Errorf("test error"),
		},
		{
			name: "No err if empty file",
			readCloserFn: func(t *testing.T, tasks []mocks2.TaskMock) io.ReadCloser {
				t.Helper()

				result := mocks2.NewMockReadCloser(ctrl)
				result.EXPECT().Read(gomock.Any()).Return(0, io.EOF).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			expectedTasks: []mocks2.TaskMock{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := getFromFile[[]mocks2.TaskMock](tt.readCloserFn(t, tt.expectedTasks))

			if tt.expectedErr != nil {
				assert.EqualError(err, tt.expectedErr.Error())
			} else {
				assert.EqualValues(tt.expectedTasks, tasks)
				assert.NoError(err)
			}
		})
	}
}
