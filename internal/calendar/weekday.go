package calendar

import "time"

const (
	sunday = time.Saturday + 1
	day    = time.Hour * 24
)

func NextSunday() *Date {
	daysUntil := sunday - time.Now().Weekday()
	nextSunday := time.Now().Add((time.Hour * 24) * time.Duration(daysUntil))

	return NewDate(time.Date(nextSunday.Year(), nextSunday.Month(), nextSunday.Day(), 0, 0, 0, 0, time.Local))
}

func NextMonday() *Date {
	return NextSunday().Add(day * time.Duration(time.Monday))
}

func NextTuesday() *Date {
	return NextSunday().Add(day * time.Duration(time.Tuesday))
}

func NextWednesday() *Date {
	return NextSunday().Add(day * time.Duration(time.Wednesday))
}

func NextThursday() *Date {
	return NextSunday().Add(day * time.Duration(time.Thursday))
}

func NextFriday() *Date {
	return NextSunday().Add(day * time.Duration(time.Friday))
}

func NextSaturday() *Date {
	return NextSunday().Add(day * time.Duration(time.Saturday))
}
