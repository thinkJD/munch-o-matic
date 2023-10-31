package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var dishId int

var orderMenu = &cobra.Command{
	Use:   "order",
	Short: "order menu",
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.OrderMenu(dishId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("enjoy")
	},
}

func init() {
	orderMenu.Flags().IntVarP(&dishId, "dish", "d", 0, "Dish to order")

	menuCmd.AddCommand(orderMenu)
}
