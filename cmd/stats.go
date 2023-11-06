package cmd

import "github.com/spf13/cobra"

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Statistics and info",
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
