package analyzer

import (
	"time"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
)

type DateParser struct {
}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (obj *DateParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		midnight         time.Time
		err              error
		moment           string
		when             *calendar.Date
		rest             []string
		isEmptyStatement = len(words) == 0
		maybeWithTime    = len(words) > 1
	)

	if isEmptyStatement {
		return words, nil
	}

	moment = words[0]
	rest = []string{}

	if midnight, err = time.Parse(calendar.DayMonthFormat, moment); err != nil {
		return words[1:], err
	}
	midnight = time.Date(time.Now().Year(), midnight.Month(), midnight.Day(), 9, 0, 0, 0, time.Local)
	when = calendar.NewDate(midnight)

	if maybeWithTime {
		if when, rest, err = when.ApplyTime(words[1:]); err != nil {
			if !errors.IsNotATimeError(err) {
				return words[1:], err
			}
		}
	}

	if when.IsToday() && when.IsPast() {
		// wrong behaviour when after 11 p.m.
		when = calendar.GetNextPeriod(calendar.NewDate(time.Now()))
	}

	if when.IsPast() {
		when = when.Add(calendar.Year)
	}

	reminder.When = when.Until()

	return rest, nil
}

func (obj *DateParser) Glossary() []string {
	var (
		startDate   = time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		finishDate  = time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
		patternList []string
	)
	for startDate.Before(finishDate) {
		patternList = append(patternList, startDate.Format(calendar.DayMonthFormat))
		startDate = startDate.Add(calendar.Day)
	}

	return patternList
}

func (obj *DateParser) PatternList() []string {
	return []string{
		`28.04 в 16:30 встреча коллектива`,
	}
}
