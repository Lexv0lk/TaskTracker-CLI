/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your task list",
	Long: `Add a new task with a short description. 
The task will be saved to the JSON storage and assigned a unique ID.

Example usage:
  task-cli add "Buy groceries"
Output:
  Task added successfully (ID: 1)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Error: Task description is required.")
			return
		}

		description := args[0]
		task, err := tasks.AddTask(description)

		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
		} else {
			cmd.Printf("Task added successfully (ID: %d)\n", task.Id)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
