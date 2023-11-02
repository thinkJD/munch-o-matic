package cmd

import (
	"fmt"
	"munch-o-matic/client"

	"github.com/spf13/cobra"
)

var countDish = &cobra.Command{
	Use:   "count",
	Short: "How often a DishId was booked",
	Run: func(cmd *cobra.Command, args []string) {
		userResp, err := cli.GetUser()
		if err != nil {
			fmt.Errorf("Error %v", err)
		}

		dishCount, dish, err := client.GetOrderCount(userResp.User.Customer.Bookings, dishId)
		if err != nil {
			fmt.Errorf("Error calculating dish count: %w", err)
		}

		fmt.Printf("%v was ordered %v times.\n", dish.Name, dishCount)

	},
}

func init() {
	countDish.Flags().IntVar(&dishId, "dish", 0, "DishId")
	menuCmd.AddCommand(countDish)
}
