package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getMenu = &cobra.Command{
	Use:   "list",
	Short: "list all menus",
	Run: func(cmd *cobra.Command, args []string) {
		if menuDay != "" && menuWeeks != 0 {
			menuWeeks = 0
			fmt.Println("--day and --week are mutually exclusive. --week is ignored")
		}

		if menuDay != "" {
			// get menu for a day
			fmt.Print("Implement me")
		} else {
			// Get menu for a week
		}

		menus, err := cli.GetMenu(menuWeeks)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "OrderId", "Booked", "Name", "Description", "Orders", "DishId"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)

		var dates []string
		for date := range menus {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		for _, date := range dates {
			for _, menu := range menus[date] {
				table.Append([]string{
					menu.Date.Format("Mon 02.01.06"),
					fmt.Sprintf("%d", menu.OrderId),
					getBookedIndicator(menu.Booked),
					menu.Dish.Name,
					menu.Dish.Description,
					fmt.Sprintf("%v", menu.Orders),
					fmt.Sprintf("%d", menu.Dish.ID)})
			}
		}
		table.Render()
	},
}

// Add any command-specific flags or arguments here

func init() {
	getMenu.Flags().IntVarP(&menuWeeks, "weeks", "w", 0, "Get Menu for n weeks")
	getMenu.Flags().StringVarP(&menuDay, "day", "d", "", "Get Menu for this day. Format: 01-02-23")
	menuCmd.AddCommand(getMenu)
}

func getBookedIndicator(b bool) string {
	if b {
		return "😋"
	}
	return " "
}
