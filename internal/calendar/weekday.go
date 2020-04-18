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
	var today = time.Now().Weekday()
	if today < time.Monday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Monday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Monday))
}

func NextTuesday() *Date {
	var today = time.Now().Weekday()
	if today < time.Tuesday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Tuesday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Tuesday))
}

func NextWednesday() *Date {
	var today = time.Now().Weekday()
	if today < time.Wednesday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Wednesday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Wednesday))
}

func NextThursday() *Date {
	var today = time.Now().Weekday()
	if today < time.Thursday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Thursday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Thursday))
}

func NextFriday() *Date {
	var today = time.Now().Weekday()
	if today < time.Friday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Friday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Friday))
}

func NextSaturday() *Date {
	var today = time.Now().Weekday()
	if today < time.Saturday {
		return NewDate(LastMidnight()).Add(day * time.Duration(time.Saturday-today))
	}

	return NextSunday().Add(day * time.Duration(time.Saturday))
}
