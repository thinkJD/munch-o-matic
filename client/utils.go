package client

import (
	"time"
)

type MenuDates struct {
	Year         int
	CalendarWeek int
}

func getNextFourWeeks() []MenuDates {
	var weeks []MenuDates
	currentTime := time.Now()

	for i := 0; i < 4; i++ {
		year, week := currentTime.AddDate(0, 0, i*7).ISOWeek()
		weeks = append(weeks, MenuDates{
			Year:         year,
			CalendarWeek: week,
		})
	}

	return weeks
}
