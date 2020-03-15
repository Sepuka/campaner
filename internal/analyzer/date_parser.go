package analyzer

import (
	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"
)

const (
	today    = `сегодня`
	tomorrow = `завтра`
)

type DateParser struct {
}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (obj *DateParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		offset = 1
		rest   []string
		word   string
		err    error
		when   *calendar.Date
	)

	if len(words) == 0 {
		return words, nil
	}
	word = words[0]

	switch word {
	case tomorrow:
		when = calendar.NewDate(calendar.NextMidnight())
	case today:
		when = calendar.NewDate(calendar.LastMidnight())
	default:
		return words, errors.NewUnConsistentGlossaryError(word, obj.Glossary())
	}

	if when, rest, err = when.ApplyTime(words[offset:]); err != nil {
		if errors.GetType(err) == errors.NotATimeError {
			switch word {
			case today:
				when = calendar.GetNextPeriod(when)
			case tomorrow:
				when = calendar.NewDate(calendar.NextMorning())
			}

		}
	}
	reminder.When = when.Until()

	return rest, nil
}

func (obj *DateParser) Glossary() []string {
	return []string{
		today,
		tomorrow,
	}
}

func (obj *DateParser) PatternList() []string {
	return []string{
		`утром позавтракать`,
		`завтра позвонить маме`,
		`завтра вечером вынести мусор`,
		`завтра в 15:35 назначить встречу`,
	}
}
