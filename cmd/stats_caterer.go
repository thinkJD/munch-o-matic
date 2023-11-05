package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var catererStats = &cobra.Command{
	Use:   "caterer",
	Short: "Some caterer related statistics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Implement Me")
		
	},
}

// Add any command-specific flags or arguments here

func init() {
	statsCmd.AddCommand(catererStats)
}
