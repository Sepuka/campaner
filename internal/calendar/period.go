package calendar

import "time"

func GetNextPeriod(date *Date) *Date {
	var hour = date.date.Hour()
	switch {
	case hour == 23:
		return NewDate(date.date.Add(10 * time.Hour))
	case hour > 17:
		return NewDate(date.date.Add(time.Duration(23-hour) * time.Hour))
	case hour > 11:
		return NewDate(date.date.Add(time.Duration(18-hour) * time.Hour))
	case hour > 8:
		return NewDate(date.date.Add(time.Duration(12-hour) * time.Hour))
	default:
		return NewDate(date.date.Add(time.Duration(9-hour) * time.Hour))
	}
}
