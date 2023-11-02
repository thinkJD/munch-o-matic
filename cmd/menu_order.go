package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var dishId int
var cancelOrder bool

var orderMenu = &cobra.Command{
	Use:   "order",
	Short: "order menu",
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.OrderMenu(dishId, cancelOrder)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("enjoy")
	},
}

func init() {
	orderMenu.Flags().IntVarP(&dishId, "dish", "d", 0, "Dish id to order")
	orderMenu.Flags().BoolVar(&cancelOrder, "c", false, "Cancel the order ")

	menuCmd.AddCommand(orderMenu)
}
