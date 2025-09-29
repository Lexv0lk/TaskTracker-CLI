/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"fmt"
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks"
	taskdomain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"github.com/spf13/cobra"
)

// markDoneCmd represents the markDone command
var markDoneCmd = &cobra.Command{
	Use:   "mark-done",
	Short: "Mark a task as done",
	Long: `Change the status of a task to "done".
This helps you keep track of completed tasks.

Example usage:
  task-cli mark-done 1`,
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

		err = tasks.UpdateTaskStatus(id, taskdomain.Done)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
		} else {
			cmd.Println("Task marked as done successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(markDoneCmd)
}
