package analyzer

import (
	"strconv"
	"time"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/speeches"
)

type DateMonthParser struct {
}

func NewDateMonthParser() *DateMonthParser {
	return &DateMonthParser{}
}

func (obj *DateMonthParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	const (
		patternLength = 2
	)

	var (
		err             error
		pattern         *speeches.Pattern
		date, monthName string
		value           int64
		ok              bool
		month           time.Month
		when            *calendar.Date
		now             = time.Now()
		months          = map[string]time.Month{
			`января`:   1,
			`февраля`:  2,
			`марта`:    3,
			`апреля`:   4,
			`мая`:      5,
			`июня`:     6,
			`июля`:     7,
			`августа`:  8,
			`сентября`: 9,
			`октября`:  10,
			`ноября`:   11,
			`декабря`:  12,
		}
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	if err = pattern.MakeOut(&date, &monthName); err != nil {
		return err
	}

	if value, err = strconv.ParseInt(date, 0, 8); err != nil {
		return errors.NewNotATimeError()
	}
	if value < calendar.MinDayOfMonth || value > calendar.MaxDayOfMonth {
		return errors.NewNotATimeError()
	}

	if month, ok = months[monthName]; !ok {
		return errors.NewNotATimeError()
	}

	if err = speech.ApplyPattern(pattern); err != nil {
		return err
	}

	when = calendar.NewDate(time.Date(now.Year(), month, int(value), 9, 0, 0, 0, time.Local))

	if when.IsPast() {
		when = when.NextYear()
	}

	if when, err = when.ApplyTime(speech); err != nil {
		if !errors.IsNotATimeError(err) {
			return err
		}
	}

	reminder.When = when.Until()

	return nil
}

func (obj *DateMonthParser) Glossary() []string {
	var (
		day  int
		days = make([]string, 31)
	)
	for day = 1; day < 32; day++ {
		days[day-1] = strconv.Itoa(day)
	}

	return days
}

func (obj *DateMonthParser) PatternList() []string {
	return []string{
		`8 января день рождения у Васеньки`,
	}
}
