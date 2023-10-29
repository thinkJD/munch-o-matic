package cmd

import (
	"github.com/spf13/cobra"
)

var getUser = &cobra.Command{
	Use:   "get-user",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetUser()
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(getUser)
}
