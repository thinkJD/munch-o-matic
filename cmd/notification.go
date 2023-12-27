package cmd

import (
	"munch-o-matic/core"

	"github.com/spf13/cobra"
)

var notification = &cobra.Command{
	Use:   "send-notification",
	Short: "Send notifications",
	Run: func(cmd *cobra.Command, args []string) {

		template := "Hello, your balance: {{.User.Customer.AccountBalance.Amount}}"

		user, err := cli.GetUser()
		if err != nil {
			print("Error: ", err)
		}

		err = core.SendTemplateNotification("thinkjd_munch_o_matic", template, user)
		if err != nil {
			print(err)
		}
	},
}

func init() {

	rootCmd.AddCommand(notification)
}
