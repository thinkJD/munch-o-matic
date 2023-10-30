package cmd

import (
	"github.com/spf13/cobra"
)

var getMenue = &cobra.Command{
	Use:   "get-menue",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetMenue()
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(getMenue)
}
