package types

import "time"

type UpcomingDish struct {
	Dummy   bool
	Date    time.Time
	OrderId int
	Dish    Dish
	Orders  int // We get the total order for each dish from the API ;-)
	Booked  bool
}
