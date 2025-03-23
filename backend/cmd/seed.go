/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
