package client

import (
	"fmt"
	"time"
)

type MenuDates struct {
	Year         int
	CalendarWeek int
}

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

// Handle a API inconsistency
// Sometimes the date is an int64 sometimes a string
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
