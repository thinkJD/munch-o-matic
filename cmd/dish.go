package cmd

import "github.com/spf13/cobra"

var dishCmdDishId int
var dishCmdOrderId int

var dishCmd = &cobra.Command{
	Use:   "dish",
	Short: "Operations related to dishes",
}

func init() {
	rootCmd.AddCommand(dishCmd)
}
