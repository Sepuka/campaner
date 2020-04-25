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
		date   time.Time
		err    error
		moment string
		when   *calendar.Date
		rest   []string
	)

	if len(words) == 0 {
		return words, nil
	}

	moment = words[0]

	if date, err = time.Parse(`01.02`, moment); err != nil {
		return words[1:], err
	}
	date = time.Date(time.Now().Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	when = calendar.NewDate(date)

	if len(words) > 1 {
		if when, rest, err = when.ApplyTime(words[1:]); err != nil {
			if errors.IsNotATimeError(err) {
				when = calendar.GetNextPeriod(when)
			}
		}
	} else {
		when = calendar.GetNextPeriod(when)
		rest = []string{}
	}

	if when.IsToday() && when.IsPast() {
		// wrong behaviour when after 11 p.m.
		when = calendar.GetNextPeriod(when)
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
		patternList = append(patternList, startDate.Format(`01.02`))
		startDate = startDate.Add(calendar.Day)
	}

	return patternList
}

func (obj *DateParser) PatternList() []string {
	return []string{
		`28.04 в 16:30 встреча коллектива`,
	}
}
