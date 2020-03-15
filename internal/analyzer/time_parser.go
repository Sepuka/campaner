package analyzer

import (
	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"
)

type TimeParser struct {
}

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (obj *TimeParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		rest []string
		err  error
		when = calendar.NewDate(calendar.LastMidnight())
	)

	if when, rest, err = when.ApplyTime(words); err != nil {
		if errors.GetType(err) == errors.NotATimeError {
			return words[1:], nil
		}
		return words, err
	}

	reminder.When = when.Until()

	return rest, err
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
