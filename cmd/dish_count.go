package cmd

import (
	"fmt"

	. "munch-o-matic/client/utils"

	"github.com/spf13/cobra"
)

var countDish = &cobra.Command{
	Use:   "count",
	Short: "How often a DishId was booked",
	Run: func(cmd *cobra.Command, args []string) {
		userResp, err := cli.GetUser()
		if err != nil {
			fmt.Println(err)
		}

		dishCount, dish, err := GetOrderCount(userResp.User.Customer.Bookings, dishId)
		if err != nil {
			fmt.Printf("Error calculating dish count: %v\n", err)
		} else {
			fmt.Printf("%v was ordered %v times.\n", dish.Name, dishCount)
		}
	},
}

func init() {
	countDish.Flags().IntVar(&dishId, "dish", 0, "DishId")
	dishCmd.AddCommand(countDish)
}