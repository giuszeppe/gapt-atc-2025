package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Manage the database for the ATC Communication Simulator",
	Long: `The database command provides utilities for managing the database
inside the ATC Communication Simulator. It includes operations for seeding,
resetting, and refreshing the database.

Available subcommands:
  - seed    → Populate the database with initial data.
  - reset   → Drop and recreate the database schema.
  - refresh → Reset and reseed the database.

Examples:
  # Seed the database
  app database seed

  # Reset the database schema
  app database reset

  # Refresh the database (reset + seed)
  app database refresh

These commands are useful for setting up, testing, and maintaining database consistency.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Manage your database with subcommands: seed, reset, refresh")
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
