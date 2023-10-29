package cmd

import (
	"fmt"

	"munch-o-matic/client"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to place an order here
		fmt.Println("Logging in...")
		// Initialize RestClient
		restClient := &client.RestClient{}

		// Perform login
		err := restClient.Login("FyGeorge2", "fmu7VJM9heb!fhg8pck")
		if err != nil {
			fmt.Println("Failed to login:", err)
			return
		}

	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(loginCmd)
}
