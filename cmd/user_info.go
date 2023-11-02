package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var userInfo = &cobra.Command{
	Use:   "info",
	Short: "TasteNext account login",
	Run: func(cmd *cobra.Command, args []string) {
		userResponse, err := cli.GetUser()
		if err != nil {
			log.Fatal("Could not load user object")
		}

		fmt.Printf("User ID:\t%v\n", userResponse.User.ID)
		fmt.Printf("Customer ID: \t%v\n", cli.CustomerId)
		fmt.Printf("Login Name:\t%v\n", userResponse.User.Username)
		fmt.Printf("First Name:\t%v\n", userResponse.User.FirstName)
		fmt.Printf("Last name:\t%v\n", userResponse.User.LastName)
		fmt.Printf("Locked:\t\t%v\n", userResponse.User.Locked)
		fmt.Println("----------------------------------")
		fmt.Printf("Account balance:\t\t%v\n", userResponse.User.Customer.AccountBalance.Amount)
		fmt.Printf("Total bookings:\t\t%v\n", len(userResponse.User.Customer.Bookings))
		totalPayed := 0
		for _, i := range userResponse.User.Customer.Bookings {
			totalPayed += i.BookingPrice
		}
		fmt.Printf("Total Payed:\t\t%v€\n", totalPayed/100)
	},
}

// Add any command-specific flags or arguments here

func init() {
	userCmd.AddCommand(userInfo)
}
