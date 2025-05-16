package cmd

import (
	"fmt"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/db"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Populate the database with initial data",
	Long: `The seed command inserts predefined data into the database.
This is useful for setting up default values, testing, or initializing an empty database.

Examples:
  # Seed the database with default data
  app database seed

  # Seed the database with test data
  app database seed --test

By default, the command seeds essential records. Use flags to customize the behavior.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database seeding started...")
		db.SeedDb()
		fmt.Println("Database seeding completed!")
	},
}

func init() {
	databaseCmd.AddCommand(seedCmd)
}
