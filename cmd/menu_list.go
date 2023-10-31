package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getMenu = &cobra.Command{
	Use:   "list",
	Short: "list all menus",
	Run: func(cmd *cobra.Command, args []string) {
		menus, err := cli.GetMenu()
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Id", "Name", "Description"})
		for _, menu := range menus {
			table.Append([]string{menu.Date.Format("Mon 02.01.06"),
				fmt.Sprintf("%d", menu.OrderId),
				menu.Dish.Name, menu.Dish.Description},
			)
		}

		table.Render()
	},
}

// Add any command-specific flags or arguments here

func init() {
	menuCmd.AddCommand(getMenu)
}
