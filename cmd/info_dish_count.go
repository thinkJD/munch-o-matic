package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoDishCount = &cobra.Command{
	Use:   "count",
	Short: "Counts how often a dish was booked in the past.",
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
	infoDishCount.Flags().IntVarP(&dishId, "dish-id", "d", 0, "ID of the dish to count")
	infoCmd.AddCommand(infoDishCount)
}
