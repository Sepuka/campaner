package calendar

import "time"

const (
	notSoon = Day * 2
)

func LastMidnight() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func NextMidnight() time.Time {
	return LastMidnight().Add(Day)
}

func NextMorning() time.Time {
	return nextPartOfADay(9)
}

func NextAfternoon() time.Time {
	return nextPartOfADay(12)
}

func NextEvening() time.Time {
	return nextPartOfADay(18)
}

func NextNight() time.Time {
	return nextPartOfADay(23)
}

func nextPartOfADay(hour int) time.Time {
	var (
		now         = time.Now()
		nextMorning = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local)
	)
	if now.Hour() >= hour {
		nextMorning = nextMorning.Add(Day)
	}

	return nextMorning
}

func IsNotSoon(when time.Duration) bool {
	return when > notSoon
}
