package calendar

import "time"

func GetNextPeriod(date *Date) *Date {
	var (
		hour       = date.date.Hour()
		nextPeriod *Date
	)

	switch {
	case hour == 23:
		nextPeriod = NewDate(date.date.Add(10 * time.Hour))
	case hour > 17:
		nextPeriod = NewDate(date.date.Add(time.Duration(23-hour) * time.Hour))
	case hour > 11:
		nextPeriod = NewDate(date.date.Add(time.Duration(18-hour) * time.Hour))
	case hour > 8:
		nextPeriod = NewDate(date.date.Add(time.Duration(12-hour) * time.Hour))
	default:
		nextPeriod = NewDate(date.date.Add(time.Duration(9-hour) * time.Hour))
	}

	return NewDate(time.Date(nextPeriod.date.Year(), nextPeriod.date.Month(), nextPeriod.date.Day(), nextPeriod.date.Hour(), 0, 0, 0, time.Local))
}
