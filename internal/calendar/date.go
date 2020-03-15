package calendar

import "time"

func LastMidnight() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func NextMidnight() time.Time {
	return LastMidnight().Add(time.Hour * 24)
}

func NextMorning() time.Time {
	return NextMidnight().Add(time.Hour * 9)
}
