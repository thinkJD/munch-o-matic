package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var dishId int
var cancelOrder bool

var orderMenu = &cobra.Command{
	Use:   "order",
	Short: "Order or cancel a dish from the menu",
	Run: func(cmd *cobra.Command, args []string) {
		if dishId == 0 {
			fmt.Println("Error: The --dish flag is mandatory")
			cmd.Usage()
			os.Exit(1)
		}
		err := cli.OrderMenu(dishId, cancelOrder)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("enjoy")
	},
}

func init() {
	orderMenu.Flags().IntVarP(&dishId, "dish", "d", 0, "DishId to order")
	orderMenu.Flags().BoolVar(&cancelOrder, "c", false, "Cancel the dish")

	menuCmd.AddCommand(orderMenu)
}
