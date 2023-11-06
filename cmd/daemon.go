package cmd

import (
	"fmt"
	cliUtils "munch-o-matic/core"

	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run munch-o-matic in daemon mode",
	Run: func(cmd *cobra.Command, args []string) {
		err := cliUtils.ValidateConfig(cfg.Daemon)
		if err != nil {
			fmt.Println(err)
		}

		err = cliUtils.Run(cfg)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
