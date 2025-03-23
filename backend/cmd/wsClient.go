/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coder/websocket"
	"github.com/spf13/cobra"
)

// wsClientCmd represents the wsClient command
var wsClientCmd = &cobra.Command{
	Use:   "wsClient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, "ws://localhost:5555/ws", nil)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "Client closing connection")

		fmt.Println("Connected to the echo server. Type messages to send (Ctrl+C to exit):")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()

			// Send message to the server
			err = conn.Write(ctx, websocket.MessageText, []byte(message))
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}

			// Read echoed message from the server
			_, data, err := conn.Read(ctx)
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}

			fmt.Printf("Echoed: %s\n", data)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from input: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(wsClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
