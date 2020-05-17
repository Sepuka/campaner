package calendar

import (
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/errors"
)

type Date struct {
	date time.Time
}

func NewDate(date time.Time) *Date {
	return &Date{
		date: date,
	}
}

func (d *Date) IsPast() bool {
	return d.date.Before(time.Now())
}

func (d *Date) IsItToday() bool {
	var now = time.Now()
	return now.Year() == d.date.Year() &&
		now.Month() == d.date.Month() &&
		now.Day() == d.date.Day()
}

func (d *Date) Add(t time.Duration) *Date {
	return NewDate(d.date.Add(t))
}

func (d *Date) Until() time.Duration {
	return time.Until(d.date)
}

func (d *Date) ApplyTime(speech *speeches.Speech) (*Date, error) {
	var (
		dateTime *Date
		err      error
	)

	if dateTime, err = findTimeOfADay(d, speech); err == nil {
		return dateTime, nil
	}

	if dateTime, err = onTimeParser(d, speech); err == nil {
		return dateTime, nil
	}

	if dateTime, err = overTimeParser(d, speech); err == nil {
		return dateTime, nil
	}

	if dateTime, err = onNextTimePeriod(d, speech); err == nil {
		return dateTime, nil
	}

	if dateTime, err = onNextTimeHourPeriod(d, speech); err == nil {
		return dateTime, nil
	}

	return d, errors.NewNotATimeError()
}

func (d *Date) Morning() *Date {
	return NewDate(d.setTime(9))
}

func (d *Date) Afternoon() *Date {
	return NewDate(d.setTime(12))
}

func (d *Date) Evening() *Date {
	return NewDate(d.setTime(18))
}

func (d *Date) Night() *Date {
	return NewDate(d.setTime(23))
}

func (d *Date) setTime(value int) time.Time {
	return d.date.Add(time.Duration(value) * time.Hour)
}
