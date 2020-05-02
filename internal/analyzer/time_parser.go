package analyzer

import (
	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/errors"
	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"
)

type TimeParser struct {
}

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (obj *TimeParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	const patternLength = 1
	var (
		err     error
		when    = calendar.NewDate(calendar.LastMidnight())
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	if when, err = when.ApplyTime(speech); err != nil {
		if errors.IsNotATimeError(err) {
			return speech.ApplyPattern(pattern)
		}
		return err
	}

	reminder.When = when.Until()

	return nil
}

func (obj *TimeParser) Glossary() []string {
	return []string{
		`через`,
		`в`,
	}
}

func (obj *TimeParser) PatternList() []string {
	return []string{
		`через минуту попить воды`,
		`через 48 секунд позвонить маме`,
		`в 13:45 сходить на обед`,
	}
}
