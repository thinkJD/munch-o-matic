package cmd

import (
	"fmt"
	"log"
	"munch-o-matic/client"
	"time"

	"github.com/spf13/cobra"
)

var getMenu = &cobra.Command{
	Use:   "list",
	Short: "list all menus",
	Run: func(cmd *cobra.Command, args []string) {
		if len(menuDay) > 0 && menuWeeks != 0 {
			menuWeeks = 0
			fmt.Println("--day and --weeks are mutually exclusive. --week is ignored")
		}

		var upcomingDishes = map[string][]client.UpcomingDish{}
		var err error
		if len(menuDay) != 0 {
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

		renderUpcomingDishesTable(upcomingDishes)
	},
}

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
