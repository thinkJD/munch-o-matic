package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var getUser = &cobra.Command{
	Use:   "get-user",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		userResponse, err := cli.GetUser()
		if err != nil {
			log.Fatal("Could not load user object")
		}

		jsonData, err := json.MarshalIndent(userResponse, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
	},
}

// Add any command-specific flags or arguments here

func init() {
	rootCmd.AddCommand(getUser)
}
