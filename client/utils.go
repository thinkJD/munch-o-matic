package client

import (
	"fmt"
	"time"
)

// Calculate the calendar weeks from now to count
func GetNextCalenderWeeks(count int) []MenuDates {
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
