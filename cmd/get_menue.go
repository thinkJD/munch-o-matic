package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var getMenue = &cobra.Command{
	Use:   "get-menue",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		menues, err := cli.GetMenue()
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err := json.MarshalIndent(menues, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(getMenue)
}
