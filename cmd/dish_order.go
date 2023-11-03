package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var orderDish = &cobra.Command{
	Use:   "order",
	Short: "Order or cancel a dish from the menu",
	Run: func(cmd *cobra.Command, args []string) {
		if dishCmdOrderId == 0 {
			fmt.Println("Error: The --order-id flag is mandatory")
			cmd.Usage()
			os.Exit(1)
		}
		err := cli.OrderMenu(orderId, cancel)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	orderDish.Flags().IntVarP(&orderId, "order-id", "o", 0, "OrderId of the dish")
	orderDish.Flags().BoolVar(&cancel, "c", false, "Cancel the order")

	dishCmd.AddCommand(orderDish)
}
