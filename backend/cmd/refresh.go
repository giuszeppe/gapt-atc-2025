package cmd

import (
	"fmt"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/db"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Reset and reseed the database",
	Long: `The refresh command drops all tables, recreates them, and then seeds the database with initial data.
This is useful for resetting the database during development or testing.

Examples:
  # Refresh the database by dropping and reseeding data
  app database refresh

This command is destructive, as it will remove all existing data. Use with caution.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database refresh started...")
		db.Refresh()
		fmt.Println("Database refresh completed!")
	},
}

func init() {
	databaseCmd.AddCommand(refreshCmd)
}
