package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all booked menues",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to place an order here
		fmt.Println("Makaroni")
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(listCmd)
}
