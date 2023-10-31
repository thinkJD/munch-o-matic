package cmd

import "github.com/spf13/cobra"

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Operations related to menus",
}

func init() {
	rootCmd.AddCommand(menuCmd)
}
