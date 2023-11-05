package cmd

import "github.com/spf13/cobra"

var dishCmd = &cobra.Command{
	Use:   "dish",
	Short: "Operations related to dishes",
}

func init() {
	rootCmd.AddCommand(dishCmd)
}
