package core

import (
	"fmt"
	"munch-o-matic/client"
)

// How often a dish was ordered in the past
func GetOrderCount(Bookings []client.Bookings, DishId int) (count int, dish client.Dish, error error) {
	for _, booking := range Bookings {
		if DishId == booking.MenuBlockLineEntry.Dish.ID {
			count++
			dish = booking.MenuBlockLineEntry.Dish
		}
	}

	if count > 0 {
		return count, dish, nil
	}

	return 0, client.Dish{}, fmt.Errorf("dish-id not found in orders")
}
