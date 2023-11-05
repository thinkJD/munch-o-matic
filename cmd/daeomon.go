package cmd

import (
	"fmt"
	cliUtils "munch-o-matic/client/utils"

	"github.com/spf13/cobra"
)

var daeomonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run munch-o-matic in daemon mode",
	Run: func(cmd *cobra.Command, args []string) {
		err := cliUtils.ValidateConfig(cfg.Daemon)
		if err != nil {
			fmt.Println(err)
		}

		cliUtils.Run(cfg)
	},
}

func init() {
	rootCmd.AddCommand(daeomonCmd)
}
