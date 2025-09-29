/*
Copyright © 2025 LEXVOLK
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "TaskTracker",
	Short: "A simple command-line task manager that stores tasks in a JSON file.",
	Long: `This is a command-line task management application designed to help users organize and track their tasks efficiently. 
The application stores all tasks in a JSON file, making it lightweight and easy to use without requiring a database.
Users can perform essential task operations directly from the command line, including adding new tasks, updating existing ones, deleting tasks, and changing task statuses to “in progress” or “done.” The application also provides flexible listing options, allowing users to view all tasks, only completed tasks, only pending tasks, or tasks currently in progress.
With simple arguments and commands, this CLI tool is perfect for anyone who wants a fast, lightweight, and easy-to-use task manager without leaving the terminal.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.TaskTracker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
