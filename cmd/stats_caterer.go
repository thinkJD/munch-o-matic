package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var catererStats = &cobra.Command{
	Use:   "caterer",
	Short: "Some user related statistics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Implement Me")

	},
}

// Add any command-specific flags or arguments here

func init() {
	statsCmd.AddCommand(catererStats)
}
