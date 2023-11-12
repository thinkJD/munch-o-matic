package cmd

import (
	"encoding/json"
	"fmt"
	"munch-o-matic/client"
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
	menuCmd.PersistentFlags().StringVarP(&formatOutput, "format", "f", "table", "Output format could be json or table.")
	rootCmd.AddCommand(menuCmd)
}

func renderOutput(DishMap client.UpcomingDishMap) error {
	switch formatOutput {
	case "table":
		renderUpcomingDishesTable(DishMap)
	case "json":
		renderUpcomingDishesJson(DishMap)
	default:
		return fmt.Errorf("%v is not a valid format. Use table or json", formatOutput)
	}
	return nil
}

func renderUpcomingDishesJson(DishMap client.UpcomingDishMap) {
	jsonData, err := json.MarshalIndent(DishMap, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling DishMap to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func renderUpcomingDishesTable(DishMap client.UpcomingDishMap) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "OrderId", "Booked", "Name", "Description", "School", "User", "DishId"})
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)

	var dates []string
	for date := range DishMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		for _, dish := range DishMap[date] {
			table.Append([]string{
				dish.Date.Format("Mon 02.01.06"),
				fmt.Sprintf("%d", dish.OrderId),
				getBookedIndicator(dish.Booked),
				dish.Dish.Name,
				dish.Dish.Description,
				fmt.Sprint(dish.Orders),
				fmt.Sprint(dish.PersonalOrders),
				fmt.Sprint(dish.Dish.ID)})
		}
	}
	table.Render()
}
