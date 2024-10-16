package memory

import (
	"errors"
	"time"
)

const invalidDateInterval = "start date must be before end date"

func daysBetween(from time.Time, to time.Time) (int, error) {
	startDate := toDay(from)
	endDate := toDay(to)

	if startDate.After(endDate) {
		return -1, errors.New(invalidDateInterval)
	}

	duration := endDate.Sub(startDate)
	days := int(duration.Hours() / 24)

	return days, nil
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		0, 0, 0, 0, time.UTC,
	)
}
