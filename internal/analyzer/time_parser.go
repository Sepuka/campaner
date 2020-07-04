package analyzer

import (
	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"
)

type TimeParser struct {
}

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (obj *TimeParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	var (
		err  error
		when = calendar.NewDate(calendar.LastMidnight())
	)

	if when, err = when.ApplyTime(speech); err != nil {
		return err
	}

	reminder.When = when.Until()

	return nil
}

func (obj *TimeParser) Glossary() []string {
	return []string{}
}

func (obj *TimeParser) PatternList() []string {
	return []string{
		`через минуту попить воды`,
		`через 48 секунд позвонить маме`,
		`в 13:45 сходить на обед`,
		`через 1 час 30 минут помыть машину`,
	}
}
