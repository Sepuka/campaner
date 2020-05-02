package analyzer

import (
	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/errors"
	"github.com/sepuka/campaner/internal/speeches"
)

const (
	morning   timesOfDay = `утром`
	afternoon timesOfDay = `днем`
	evening   timesOfDay = `вечером`
	tonight   timesOfDay = `ночью`
)

type (
	timesOfDay       string
	TimesOfDayParser struct {
	}
)

func (t timesOfDay) String() string {
	return string(t)
}

func NewTimesOfDayParser() *TimesOfDayParser {
	return &TimesOfDayParser{}
}

func (obj *TimesOfDayParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	var (
		err     error
		when    = calendar.NewDate(calendar.LastMidnight())
		pattern *speeches.Pattern
	)

	if when, err = when.ApplyTime(speech); err != nil {
		if errors.IsNotATimeError(err) {
			return speech.ApplyPattern(pattern)
		}
		return err
	}

	if when.IsPast() {
		when = when.Add(calendar.Day)
	}

	reminder.When = when.Until()

	return nil
}

func (obj *TimesOfDayParser) Glossary() []string {
	return []string{
		morning.String(),
		afternoon.String(),
		evening.String(),
		tonight.String(),
	}
}

func (obj *TimesOfDayParser) PatternList() []string {
	return []string{
		`утром совещание`,
		`днем позвонить коллеге`,
		`вечером ужин`,
	}
}
