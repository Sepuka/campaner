package analyzer

import (
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
)

type DateParser struct {
}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (obj *DateParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	const patternLength = 1
	var (
		midnight time.Time
		err      error
		pattern  *speeches.Pattern
		when     *calendar.Date
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	if midnight, err = time.Parse(calendar.DayMonthFormat, pattern.Origin()); err != nil {
		return err
	}
	midnight = time.Date(time.Now().Year(), midnight.Month(), midnight.Day(), 9, 0, 0, 0, time.Local)
	when = calendar.NewDate(midnight)

	if err = speech.ApplyPattern(pattern); err != nil {
		return err
	}

	if when, err = when.ApplyTime(speech); err != nil {
		if !errors.IsNotATimeError(err) {
			return err
		}
	}

	if when.IsItToday() && when.IsPast() {
		// wrong behaviour when after 11 p.m.
		when = calendar.GetNextPeriod(calendar.NewDate(time.Now()))
	}

	if when.IsPast() {
		when = when.Add(calendar.Year)
	}

	reminder.When = when.Until()

	return nil
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
