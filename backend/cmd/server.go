package cmd

import (
	"database/sql"
	"fmt"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/api"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "example.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		tokenStore := stores.NewTokenStore()
		userStore := stores.NewUserStore(db)
		scenarioStore := stores.NewScenarioStore(db)
		if err := api.Run(os.Getenv, tokenStore, userStore, scenarioStore); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
