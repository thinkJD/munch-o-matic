package cmd

import (
	"fmt"
	"munch-o-matic/core"

	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run munch-o-matic in daemon mode",
	Run: func(cmd *cobra.Command, args []string) {
		d, err := core.NewDaemon(cfg.Core, cli)
		if err != nil {
			fmt.Println(err)
		}

		err = d.Run()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
