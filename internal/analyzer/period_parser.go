package analyzer

import (
	"strconv"
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/errors"
	"github.com/sepuka/campaner/internal/speeches"
)

type PeriodParser struct {
}

func NewPeriodParser() *PeriodParser {
	return &PeriodParser{}
}

func (obj *PeriodParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	const (
		momentBitSize   = 32
		maxDaysPeriod   = 1000
		maxWeekPeriod   = 500
		maxMonthsPeriod = 100
		patternLength   = 3
	)
	var (
		err                          error
		pattern                      *speeches.Pattern
		preposition, size, dimension string
		value                        float64
		period                       int
		now                          = time.Now()
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	if err = pattern.MakeOut(&preposition, &size, &dimension); err != nil {
		return err
	}

	if preposition != `через` {
		return errors.NewUnexpectedPrepositionError(preposition, `через`)
	}

	if value, err = strconv.ParseFloat(size, momentBitSize); err != nil {
		return err
	}
	period = int(value)
	if period < 1 {
		return errors.NewInvalidTimeValueError(period)
	}

	switch {
	case strings.HasPrefix(dimension, `дней`), strings.HasPrefix(dimension, `дня`):
		if period > maxDaysPeriod {
			return errors.NewInvalidTimeValueError(period)
		}
		reminder.When = calendar.NewDate(now).Morning().Add(calendar.Day * time.Duration(period)).Until()
	case strings.HasPrefix(dimension, `недел`):
		if period > maxWeekPeriod {
			return errors.NewInvalidTimeValueError(period)
		}
		reminder.When = calendar.NewDate(now).Morning().Add(calendar.Day * time.Duration(period*7)).Until()
	case strings.HasPrefix(dimension, `месяц`):
		if period > maxMonthsPeriod {
			return errors.NewInvalidTimeValueError(period)
		}
		reminder.When = calendar.NewDate(now).Morning().Add(calendar.Day * time.Duration(period*30)).Until()
	}

	return nil
}

func (obj *PeriodParser) Glossary() []string {
	return []string{}
}

func (obj *PeriodParser) PatternList() []string {
	return []string{
		`через 5 дней годовщина`,
		`через 10 лет я стану на 10 лет старше`,
	}
}
