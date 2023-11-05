package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var daeomonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Operations related to dishes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cfg.Daemon)
	},
}

func init() {
	rootCmd.AddCommand(daeomonCmd)
}
