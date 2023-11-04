package cmd

import (
	"fmt"
	"log"
	"munch-o-matic/client"
	"os"
	"sort"
	"time"

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

		var upcomingDishes = map[string][]client.UpcomingDish{}
		var err error
		if menuDay != "" {
			parsedDate, err := time.Parse("02.01.06", menuDay)
			if err != nil {
				log.Fatalf("Invalid date format. Please provide the date in the format 02.01.06")
			}
			upcomingDishes, err = cli.GetMenuDay(parsedDate)
			if err != nil {
				log.Fatal("Error getting dishes: %w", err)
			}
		} else if menuWeeks != 0 {
			upcomingDishes, err = cli.GetMenuWeeks(menuWeeks)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Please provide --day or --weeks")
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "OrderId", "Booked", "Name", "Description", "Orders", "DishId"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)

		var dates []string
		for date := range upcomingDishes {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		for _, date := range dates {
			for _, dish := range upcomingDishes[date] {
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
