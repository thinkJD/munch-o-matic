package cmd

import (
	"munch-o-matic/core"

	"github.com/spf13/cobra"
)

var notification = &cobra.Command{
	Use:   "send-notification",
	Short: "Send notifications",
	Run: func(cmd *cobra.Command, args []string) {
		core.SendAccountBalanceNotification(3000)
	},
}

func init() {

	rootCmd.AddCommand(notification)
}
