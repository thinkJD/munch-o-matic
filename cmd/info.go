package cmd

import "github.com/spf13/cobra"

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information and statistics",
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
