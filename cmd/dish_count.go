package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var countDish = &cobra.Command{
	Use:   "count",
	Short: "How often a DishId was booked",
	Run: func(cmd *cobra.Command, args []string) {
		dishCount, dish := cli.GetOrderCount(dishId)
		if dishCount > 0 {
			fmt.Println("Dish was ordered 0 times.")
		} else {
			fmt.Printf("Dish %v was ordered %v times.", dishCount, dish.Name)
		}
	},
}

func init() {
	countDish.Flags().IntVarP(&dishId, "dish-id", "d", 0, "DishId")
	dishCmd.AddCommand(countDish)
}
