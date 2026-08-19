package utils

import "time"

func MakeDeadline(startTime time.Time, days int) time.Time {
	loc := time.FixedZone("UTC+5", 5*60*60)

	now := time.Now().In(loc)
	target := now.AddDate(0, 0, days)

	year, month, day := target.Date()
	deadline := time.Date(year, month, day, 23, 59, 0, 0, loc)

	return deadline
}
