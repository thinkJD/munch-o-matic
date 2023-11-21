package cmd

import (
	"fmt"
	"log"
	"munch-o-matic/core"
	"time"

	"github.com/spf13/cobra"
)

var autoOrderMenu = &cobra.Command{
	Use:   "auto-order",
	Short: "Order dishes automatically",
	Run: func(cmd *cobra.Command, args []string) {
		if len(menuDay) > 0 && menuWeeks != 0 {
			menuWeeks = 0
			fmt.Println("--day and --weeks are mutually exclusive. --week is ignored")
		}

		var orderedDishes core.OrderedDishes
		var err error
		if len(menuDay) != 0 {
			parsedDate, err := time.Parse("02.01.06", menuDay)
			if err != nil {
				log.Fatalf("Invalid date format. Please provide the date in the format 02.01.06")
			}
			orderedDishes, err = core.AutoOrderDay(cli, parsedDate, autoOrderStrategy, dryRun)
			if err != nil {
				log.Fatal("can't order: %w", err)
			}
		} else if menuWeeks != 0 {
			orderedDishes, err = core.AutoOrderWeek(cli, menuWeeks, 0, autoOrderStrategy, dryRun)
			if err != nil {
				log.Fatal("can't order: %w", err)
			}
		} else {
			log.Fatal("Please provide --day or --weeks")
		}
		fmt.Println(orderedDishes)
	},
}

// Add any command-specific flags or arguments here

func init() {
	autoOrderMenu.Flags().IntVarP(&menuWeeks, "weeks", "w", 0, "Get Menu for n weeks")
	autoOrderMenu.Flags().StringVarP(&menuDay, "day", "d", "", "Get Menu for this day. Format: 01-02-23")

	autoOrderMenu.Flags().StringVarP(&autoOrderStrategy, "strategy", "s", "Random", "Strategy used to pick dish")
	autoOrderMenu.Flags().BoolVar(&dryRun, "dry-run", false, "Do not order but print out dish picks")
	menuCmd.AddCommand(autoOrderMenu)
}
