package cmd

import (
	"fmt"
	. "munch-o-matic/client/types"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Operations related to menus",
}

func init() {
	rootCmd.AddCommand(menuCmd)
}

func renderUpcomingDishesTable(UpcomingDishes map[string][]UpcomingDish) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "OrderId", "Booked", "Name", "Description", "Orders", "DishId"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	var dates []string
	for date := range UpcomingDishes {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		for _, dish := range UpcomingDishes[date] {
			table.Append([]string{
				dish.Date.Format("Mon 02.01.06"),
				fmt.Sprintf("%d", dish.OrderId),
				getBookedIndicator(dish.Booked),
				dish.Dish.Name,
				dish.Dish.Description,
				fmt.Sprintf("%v", dish.Orders),
				fmt.Sprintf("%d", dish.Dish.ID)})
		}
	}
	table.Render()
}
