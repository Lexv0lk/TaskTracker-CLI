[![GitHub - lakshay1341/Task-Tracker: Task Tracker CLI is a command-line interface application designed to help you manage your tasks efficiently. This project allows you to add, update, delete, and list tasks, storing them in a JSON file for persistence.](https://images.openai.com/thumbnails/url/JDPGG3icu1mUUVJSUGylr5-al1xUWVCSmqJbkpRnoJdeXJJYkpmsl5yfq5-Zm5ieWmxfaAuUsXL0S7F0Tw4u8nMzdCpINyuqMMh0MU739owqLvKyiAh1zPXPi3QNc0zMzkkO8S1PLKnwNPMpM83zcAuOTDJ2NfdQKwYAweco4g)](https://github.com/lakshay1341/Task-Tracker?utm_source=chatgpt.com)

# TaskTracker-CLI

TaskTracker-CLI is a command-line application developed in Go to help you efficiently manage your tasks. It allows you to add, update, delete, and list tasks, storing them in a JSON file for persistence.

---

## Features

* **Add a new task**: Create a new task with a description.
* **Update an existing task**: Modify the description of an existing task.
* **Delete a task**: Remove a task from the list.
* **Mark a task as in progress**: Change the status of a task to "in-progress".
* **Mark a task as done**: Change the status of a task to "done".
* **List all tasks**: Display all tasks with their current status.
* **List tasks by status**: Filter tasks based on their status (done, todo, in-progress).
* **Task storage**: All tasks are stored locally in a JSON file in your home directory (Windows: %AppData%\Roaming\TaskTracker-CLI)
---

## Installation

Required go to be installed.

   ```bash
   go install github.com/Lexv0lk/TaskTracker-CLI/task-cli@v1.0.0
   ```

The binary will be installed in your $GOPATH/bin or $HOME/go/bin.

---

## Usage

Run the application with the desired command:

### Add a new task

```bash
task-cli add "Buy groceries"
```

### Update an existing task

```bash
task-cli update 1 "Buy groceries and cook dinner"
```

### Delete a task

```bash
task-cli delete 1
```

### Mark a task as in progress

```bash
task-cli mark-in-progress 1
```

### Mark a task as done

```bash
task-cli mark-done 1
```

### List all tasks

```bash
task-cli list
```

### List tasks by status

```bash
task-cli list done
task-cli list todo
task-cli list in-progress
```

---

## License

This project is licensed under the MIT License.

---

## Contributing

Feel free to fork the repository, submit issues, and send pull requests. Contributions are welcome!

---

## Acknowledgements

This project was developed to practice Go programming.

---

For more information and to access the source code, visit the [TaskTracker-CLI GitHub repository](https://github.com/Lexv0lk/TaskTracker-CLI).
