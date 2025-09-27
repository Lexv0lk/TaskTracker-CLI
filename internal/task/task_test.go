//go:generate mockgen -source storage.go -destination=mocks/storage.go -package=mocks
package task

import "testing"

func TestAddTask(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		addingTask    task
		testStorageFn func(t *testing.T, addingTask task) TaskStorage
		expectedErr   error
	}
}
