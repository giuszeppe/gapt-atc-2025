package cmd

import (
	"fmt"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/db"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Drop and recreate the database schema",
	Long: `The reset command removes all tables and recreates the database schema
without seeding any data. This is useful for starting from a clean slate
without any predefined records.

Examples:
  # Reset the database (drop and recreate tables)
  app database reset

⚠️ Warning: This operation is irreversible and will permanently delete all data in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database reset started...")
		db.ResetDb()
		// Add logic to drop and recreate the database schema
		fmt.Println("Database reset completed!")
	},
}

func init() {
	databaseCmd.AddCommand(resetCmd)
}
