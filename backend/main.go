/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"github.com/giuszeppe/gatp-atc-2025/backend/cmd"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal"
)

func main() {
	err := internal.LoadEnv(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	cmd.Execute()
}
