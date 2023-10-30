package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getMenue = &cobra.Command{
	Use:   "get-menue",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		menus, err := cli.GetMenue()
		if err != nil {
			log.Fatal(err)
		}

		// Initialize tablewriter
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"kw", "Name", "Description"})

		// Iterate through the list and append each dish's name and description to the table
		for _, menu := range menus {
			calendarWeek := menu.MenuBlockWeekWrapper.MenuBlockWeek.CalendarWeek
			for _, week := range menu.MenuBlockWeekWrapper.MenuBlockWeek.MenuBlockLineWeeks {
				for _, entry := range week.Entries {
					dishName := entry.Dish.Name
					dishDescription := entry.Dish.Description
					table.Append([]string{fmt.Sprintf("%d", calendarWeek), dishName, dishDescription})
				}
			}
		}

		table.Render()
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(getMenue)
}
