/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// databaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
