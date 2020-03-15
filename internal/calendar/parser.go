package calendar

import (
	"strconv"
	"time"

	"github.com/sepuka/campaner/internal/domain"

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

func (d *Date) Until() time.Duration {
	return time.Until(d.date)
}

func (d *Date) ApplyTime(words []string) (*Date, []string, error) {
	var (
		dateTime *Date
		err      error
		rest     []string
	)

	if len(words) == 0 {
		return nil, words, errors.NewNotATimeError()
	}

	if dateTime, rest, err = d.findTimeOfDay(words); err == nil {
		return dateTime, rest, nil
	}

	if dateTime, rest, err = d.onTimeParser(words); err == nil {
		return dateTime, rest, nil
	}

	if dateTime, rest, err = d.overTimeParser(words); err == nil {
		return dateTime, rest, nil
	}

	return nil, words, errors.NewNotATimeError()
}

func (d *Date) findTimeOfDay(words []string) (*Date, []string, error) {
	var moment = words[0]

	switch moment {
	case `утром`:
		return NewDate(d.setTime(9)), words[1:], nil
	case `днем`:
		return NewDate(d.setTime(12)), words[1:], nil
	case `вечером`:
		return NewDate(d.setTime(18)), words[1:], nil
	case `ночью`:
		return NewDate(d.setTime(23)), words[1:], nil
	default:
		return nil, words, errors.NewNotATimeError()
	}
}

func (d *Date) onTimeParser(words []string) (*Date, []string, error) {
	var (
		exactTime time.Time
		err       error
		moment    string
	)

	if len(words) < 2 {
		return nil, words, errors.NewNotATimeError()
	}

	if words[0] != `в` {
		return nil, words, errors.NewNotATimeError()
	}

	moment = words[1]
	if exactTime, err = time.Parse(`15:04:05`, moment); err != nil {
		if exactTime, err = time.Parse(`15:04`, moment); err != nil {
			return nil, words, errors.NewNotATimeError()
		}
	}

	exactTime = time.Date(d.date.Year(), d.date.Month(), d.date.Day(), exactTime.Hour(), exactTime.Minute(), 0, 0, time.Local)

	if exactTime.Before(time.Now()) {
		exactTime = exactTime.Add(24 * time.Hour)
	}

	return NewDate(exactTime), words[2:], nil
}

func (d *Date) overTimeParser(words []string) (*Date, []string, error) {
	var (
		value      float64
		dimension  string
		restOffset int
		err        error
		timeFrame  *domain.TimeFrame
		duration   time.Duration
		moment     string
	)

	if len(words) < 2 {
		return nil, words, errors.NewNotATimeError()
	}

	if words[0] != `через` {
		return nil, words, errors.NewNotATimeError()
	}

	moment = words[1]
	if value, err = strconv.ParseFloat(moment, 32); err != nil {
		value = 1
		dimension = moment
		restOffset = 2
	} else {
		dimension = words[2]
		restOffset = 3
	}

	timeFrame = domain.NewTimeFrame(value, dimension)
	if duration, err = timeFrame.GetDuration(); err != nil {
		return nil, words, errors.NewNotATimeError()
	}

	return NewDate(time.Now().Add(duration)), words[restOffset:], nil
}

func (d *Date) setTime(value int) time.Time {
	return d.date.Add(time.Duration(value) * time.Hour)
}
