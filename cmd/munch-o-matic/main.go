package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "munch-o-matic",
		Short: "TasteNext API Client",
	}

	// Add the order command to the root command
	rootCmd.AddCommand(Login())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
