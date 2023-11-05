package client

import (
	"fmt"
	"math/rand"
	"time"
)

type MenuDates struct {
	Year         int
	CalendarWeek int
}

// Calculate the calendar weeks from now to count
func getNextCalenderWeeks(count int) []MenuDates {
	var weeks []MenuDates
	currentTime := time.Now()

	for i := 0; i < count; i++ {
		year, week := currentTime.AddDate(0, 0, i*7).ISOWeek()
		weeks = append(weeks, MenuDates{
			Year:         year,
			CalendarWeek: week,
		})
	}
	return weeks
}

// Handle a API inconsistency. The date can be string or int64
func GetEmissionDateAsTime(emissionDate interface{}) (time.Time, error) {
	switch v := emissionDate.(type) {
	case float64:
		return time.Unix(int64(v)/1000, (int64(v)%1000)*1e6), nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return time.Time{}, fmt.Errorf("parsing time from string failed")
		}
		return t, nil
	default:
		return time.Time{}, fmt.Errorf("unknown type for mission_date: %T", v)
	}
}

// How often a dish was ordered in the past
func GetOrderCount(Bookings []Bookings, DishId int) (count int, dish Dish, error error) {
	for _, booking := range Bookings {
		if DishId == booking.MenuBlockLineEntry.Dish.ID {
			count++
			dish = booking.MenuBlockLineEntry.Dish
		}
	}

	if count > 0 {
		return count, dish, nil
	}

	return 0, Dish{}, fmt.Errorf("Dish ID not found in orders")
}

// Pick dishes automatically based on a few strategies
func ChooseDishesByStrategy(Strategy string, UpcomingDishes map[string][]UpcomingDish) (map[int]UpcomingDish, error) {
	retVal := map[int]UpcomingDish{}

	// Helper function to decide if menu should be skipped
	shouldSkipMenu := func(menu []UpcomingDish) bool {
		for _, dish := range menu {
			if dish.Booked || dish.Dummy {
				return true
			}
		}
		return false
	}

	// Iterate over the dishes of the day
	for _, menu := range UpcomingDishes {
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
			return map[int]UpcomingDish{}, fmt.Errorf("PersonalFav is not implemented, sorry")

		default:
			return map[int]UpcomingDish{}, fmt.Errorf("%v is not a valid strategy", Strategy)
		}
	}
	return retVal, nil
}
