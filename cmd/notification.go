package cmd

import (
	"log"
	"munch-o-matic/core"

	"github.com/spf13/cobra"
)

var notification = &cobra.Command{
	Use:   "send-notification",
	Short: "Send notifications",
	Run: func(cmd *cobra.Command, args []string) {

		err := core.SendNotification("thinkjd_munch_o_matic", title, message)
		if err != nil {
			log.Fatal("could not send notification: ", err)
		}
	},
}

func init() {
	notification.Flags().StringVarP(&ntfyTopic, "topic", "T", "", "Ntfy topic")
	notification.Flags().StringVarP(&title, "title", "t", "CLI", "Notification title")
	notification.Flags().StringVarP(&message, "message", "m", "", "Notificatoin message")
	rootCmd.AddCommand(notification)
}
