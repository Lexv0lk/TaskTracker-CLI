package task

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	saveFileName    = "tasks.json"
	defaultSavePath = "%AppData%"
	appName         = "TaskTracker-CLI"
)

func saveToFile(file io.WriteCloser, tasks []task) error {
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}

func getFromFile(file io.ReadCloser) ([]task, error) {
	defer file.Close()

	var tasks []task
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&tasks)
	return tasks, err
}

func ensureSaveDirExists() error {
	_, err := os.Stat(getSaveDir())

	if os.IsNotExist(err) {
		return os.MkdirAll(getSaveDir(), 0755)
	}

	return err
}

func openOrCreateSaveFile() (*os.File, error) {
	if err := ensureSaveDirExists(); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(getSavePath(), os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func getSavePath() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName, saveFileName)
}

func getSaveDir() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName)
}
