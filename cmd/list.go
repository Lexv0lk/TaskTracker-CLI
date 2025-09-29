/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"github.com/Lexv0lk/TaskTracker-CLI/internal/application/tasks"
	taskdomain "github.com/Lexv0lk/TaskTracker-CLI/internal/domain/tasks"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks, optionally filtered by status",
	Long: `Display tasks from the JSON storage. 
If no status argument is provided, all tasks are listed. 
You can optionally filter tasks by status: "done", "todo", or "in-progress".

Example usage:
  # List all tasks
  task-cli list

  # List only completed tasks
  task-cli list done

  # List only pending tasks
  task-cli list todo

  # List tasks currently in progress
  task-cli list in-progress`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			res, err := tasks.GetAllTasks()

			if err != nil {
				cmd.Printf("Error: %s\n", err.Error())
			} else {
				cmd.Print(res)
			}
		} else if len(args) == 1 {
			progressStr := strings.ToLower(args[0])
			possibleProgressStrs := map[string]interface{}{
				taskdomain.TodoStr:       struct{}{},
				taskdomain.InProgressStr: struct{}{},
				taskdomain.DoneStr:       struct{}{},
			}

			if _, ok := possibleProgressStrs[progressStr]; !ok {
				cmd.Printf("Error: Invalid progress filter. Use '%s', '%s', or '%s'.\n",
					taskdomain.TodoStr, taskdomain.InProgressStr, taskdomain.DoneStr)
				return
			} else {
				progress, err := tasks.ParseStatusString(progressStr)

				if err != nil {
					cmd.Printf("Error: %s\n", err.Error())
					return
				}

				res, err := tasks.GetTasks(progress)

				if err != nil {
					cmd.Printf("Error: %s\n", err.Error())
				} else {
					cmd.Print(res)
				}
			}

		} else {
			cmd.Println("Error: No arguments or only progress filter are required.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
