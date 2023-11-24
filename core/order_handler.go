package core

import (
	"fmt"
	"log"
	"math/rand"
	"munch-o-matic/client"
	"time"
)

type OrderedDishes map[int]client.UpcomingDish

func orderDishes(Cli *client.RestClient, Dishes OrderedDishes, DryRun bool) (OrderedDishes, error) {
	retVal := OrderedDishes{}
	for _, dish := range Dishes {
		if !DryRun {
			err := Cli.OrderDish(dish.OrderId, false)
			if err != nil {
				// Return the already ordered dishes
				return retVal, fmt.Errorf("order dish %d failed with: %w", dish.OrderId, err)
			}
			retVal[dish.OrderId] = dish
		}
	}
	return retVal, nil
}

// AutoOrderWeek places orders for dishes for a specified week and year, based on a chosen strategy.
// It supports a dry-run mode for simulation. It returns a map of ordered dishes and an error, if any.
//
// Parameters:
// - Cli: Client interface to interact with the ordering system.
// - Week: The week number for which the order is being placed.
// - Year: The year for which the order is being placed. Defaults to the current year if set to 0.
// - Strategy: Strategy for selecting dishes (e.g., "Random", "SchoolFav").
// - DryRun: If true, simulates the ordering process without actual orders being placed.
//
// Returns:
// - OrderedDishes: A map of ordered dishes with their IDs.
// - error: An error, if any occurred during the ordering process.
func AutoOrderWeek(Cli *client.RestClient, Week int, Year int, Strategy string, DryRun bool) (OrderedDishes, error) {
	if Year == 0 {
		// Defaults to current year
		Year = time.Now().Year()
	}
	if Week == 0 {
		// Defaults to current week
		Year, Week = time.Now().ISOWeek()
	}

	menu, err := Cli.GetMenuWeek(Year, Week)
	if err != nil {
		return OrderedDishes{}, fmt.Errorf("error: %w", err)
	}

	dishes, err := ChooseDishesByStrategy(Strategy, menu)
	if err != nil {
		return OrderedDishes{}, fmt.Errorf("error picking dishes: %w", err)
	}

	orderedDishes, err := orderDishes(Cli, dishes, DryRun)
	if err != nil {
		return OrderedDishes{}, fmt.Errorf("order dish: %w", err)
	}

	return orderedDishes, nil
}

// AutoOrderDay places orders for one dish for a specified day, based on a chosen strategy.
// It supports a dry-run mode for simulation. It returns a map of ordered dishes and an error, if any.
//
// Parameters:
// - Cli: Client interface to interact with the ordering system.
// - Day: The day for which the order is being placed.
// - Strategy: Strategy for selecting dishes (e.g., "SchoolFav", "Random").
// - DryRun: If true, simulates the ordering process without actual orders being placed.
//
// Returns:
// - OrderedDishes: A map of ordered dishes with their IDs.
// - error: An error, if any occurred during the ordering process.
func AutoOrderDay(Cli *client.RestClient, Day time.Time, Strategy string, DryRun bool) (OrderedDishes, error) {
	upcomingDishes, err := Cli.GetMenuDay(Day)
	if err != nil {
		log.Fatal("get dishes: %w", err)
	}

	dishes, err := ChooseDishesByStrategy(Strategy, upcomingDishes)
	if err != nil {
		return OrderedDishes{}, fmt.Errorf("picking dish: %w", err)
	}

	orderedDishes, err := orderDishes(Cli, dishes, DryRun)
	if err != nil {
		return OrderedDishes{}, fmt.Errorf("order dish: %w", err)
	}

	return orderedDishes, nil
}

// ChooseDishesByStrategy selects dishes based on a specified strategy from a given map of upcoming dishes.
// Currently supported strategies include "SchoolFav" and "Random". "PersonalFav" is planned but not yet implemented.
//
// Parameters:
// - Strategy: The strategy to use for selecting dishes.
// - DishMap: A map of upcoming dishes to choose from.
//
// Returns:
// - OrderedDishes: A map of selected dishes based on the strategy.
// - error: An error if the strategy is invalid or not implemented.
func ChooseDishesByStrategy(Strategy string, DishMap client.UpcomingDishMap) (OrderedDishes, error) {
	retVal := OrderedDishes{}

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
			return OrderedDishes{}, fmt.Errorf("PersonalFav is not implemented, sorry")

		default:
			return OrderedDishes{}, fmt.Errorf("%v is not a valid strategy", Strategy)
		}
	}
	return retVal, nil
}
