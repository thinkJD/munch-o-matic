package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Operations related to users",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
