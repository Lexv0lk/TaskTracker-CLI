/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"fmt"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the description of an existing task",
	Long: `Modify the description of a task identified by its ID.
This allows you to change the task details without creating a new task.

Example usage:
  task-cli update 1 "Buy groceries and cook dinner"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Println("Error: Task ID and new description are required.")
			return
		}

		var id int
		_, err := fmt.Sscanf(args[0], "%d", &id)
		if err != nil {
			cmd.Println("Error: Invalid task ID format.")
			return
		}

		newDescription := args[1]
		err = tasks.UpdateTask(id, newDescription)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
		} else {
			cmd.Println("Task updated successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
