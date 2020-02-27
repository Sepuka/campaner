package calendar

import (
	"errors"
	"time"
)

type Date struct {
	date time.Time
}

func NewDate(date time.Time) *Date {
	return &Date{
		date: date,
	}
}

func (d *Date) GetMorning() time.Time {
	return d.date.Add(9 * time.Hour)
}

func (d *Date) GetTime(part string) (time.Time, error) {
	switch part {
	case `утром`:
		return d.GetMorning(), nil
	case `днем`:
		return d.date.Add(12 * time.Hour), nil
	case `вечером`:
		return d.date.Add(18 * time.Hour), nil
	case `ночью`:
		return d.date.Add(23 * time.Hour), nil
	default:
		return time.Now(), errors.New(`invalid part of a day`)
	}
}
