package core

import (
	"fmt"
	"math/rand"
	"munch-o-matic/client"
	"time"
)

// AutoOrder automatically places an order for upcoming dishes starting from the specified start date.
// It returns a map of upcoming dishes and an error, if any.
//
// Parameters:
// - StartDate: The start date from which to place the order.
// - Days: The number of days for which to place the order.
//
// Returns:
// - client.UpcomingDishMap: A map of upcoming dishes.
// - error: An error, if any.
func AutoOrderWeeks(Cli client.RestClient, StartDate time.Time, Days int) (client.UpcomingDishMap, error) {
	return client.UpcomingDishMap{}, fmt.Errorf("implement me")
}

// Pick dishes automatically based on a few strategies
func ChooseDishesByStrategy(Strategy string, DishMap client.UpcomingDishMap) (map[int]client.UpcomingDish, error) {
	retVal := map[int]client.UpcomingDish{}

	// Helper function to decide if menu should be skipped
	shouldSkipMenu := func(menu []client.UpcomingDish) bool {
		for _, dish := range menu {
			if dish.Booked || dish.Dummy {
				return true
			}
		}
		return false
	}

	// Iterate over the dishes of the day
	for _, menu := range DishMap {
		if shouldSkipMenu(menu) {
			continue
		}
		// Choose dish based on the strategy
		switch Strategy {

		case "SchoolFav":
			var maxPos, maxVal int
			for i, dish := range menu {
				if dish.Orders > maxVal {
					maxPos = i
					maxVal = dish.Orders
				}
			}
			retVal[menu[maxPos].OrderId] = menu[maxPos]

		case "Random":
			randomInt := rand.Intn(len(menu))
			retVal[menu[randomInt].OrderId] = menu[randomInt]

		case "PersonalFav":
			/* TODO: Add personal order count in getMenuWeek or structure the code better.
			var maxPos, maxVal int
			for i, dish := range menu {
				GetOrderCount()
			}
			*/
			return map[int]client.UpcomingDish{}, fmt.Errorf("PersonalFav is not implemented, sorry")

		default:
			return map[int]client.UpcomingDish{}, fmt.Errorf("%v is not a valid strategy", Strategy)
		}
	}
	return retVal, nil
}
