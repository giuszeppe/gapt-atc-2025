package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "Backend service for managing application data and APIs",
	Long: `The backend service provides a command-line interface to manage
the application's database, APIs, and other server-side operations.

Examples:
  # Start the backend server
  backend server

  # Reset the database
  backend database reset`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
