package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Login() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "TasteNext account login",
		Run: func(cmd *cobra.Command, args []string) {
			// Implement the logic to place an order here
			fmt.Println("Logging in...")
		},
	}

	// Add any command-specific flags or arguments here

	return cmd
}
