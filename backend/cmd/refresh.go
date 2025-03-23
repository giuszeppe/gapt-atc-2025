/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
