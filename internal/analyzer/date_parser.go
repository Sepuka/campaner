package analyzer

import (
	"errors"
	"time"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"
)

type DateParser struct {
}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (obj *DateParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		word      string
		offset    = 1
		exactly   = `утром`
		exactTime time.Time
		err       error
		date      *calendar.Date
	)

	if len(words) == 0 {
		return words, nil
	}

	word = words[0]
	if len(words) > 1 {
		exactly = words[1]
		offset++
	}

	switch word {
	case `завтра`:
		date = calendar.NewDate(calendar.NextMidnight())
	default:
		date = calendar.NewDate(calendar.LastMidnight())
	}

	if exactTime, err = date.GetTime(exactly); err != nil {
		exactTime = date.GetMorning()
	}
	reminder.When = time.Until(exactTime)

	if !reminder.IsValid() {
		return words, errors.New(`date is not valid`)
	}

	return words[offset:], nil
}

func (obj *DateParser) Glossary() []string {
	return []string{
		`сегодня`,
		`завтра`,
	}
}

func (obj *DateParser) PatternList() []string {
	return []string{
		`утром позавтракать`,
		`завтра позвонить маме`,
		`завтра вечером вынести мусор`,
	}
}
