/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"fmt"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task by its ID",
	Long: `Remove a task permanently from the task list using its ID.
This action cannot be undone.

Example usage:
  task-cli delete 1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Println("Error: Only Task ID is required.")
			return
		}

		var id int
		_, err := fmt.Sscanf(args[0], "%d", &id)
		if err != nil {
			cmd.Println("Error: Invalid task ID format.")
			return
		}

		err = tasks.DeleteTask(id)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
		} else {
			cmd.Println("Task deleted successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
