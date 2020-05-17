package analyzer

import (
	"time"

	"github.com/sepuka/campaner/internal/errors"
	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"
)

const (
	today     dayName = `сегодня`
	tomorrow  dayName = `завтра`
	monday    dayName = `понедельник`
	tuesday   dayName = `вторник`
	wednesday dayName = `среду`
	thursday  dayName = `четверг`
	friday    dayName = `пятницу`
	saturday  dayName = `субботу`
	sunday    dayName = `воскресенье`
)

type (
	dayName   string
	DayParser struct {
	}
)

func (dn dayName) String() string {
	return string(dn)
}

func NewDayParser() *DayParser {
	return &DayParser{}
}

func (obj *DayParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	var (
		err  error
		when *calendar.Date
	)

	if when, err = obj.findDay(speech); err != nil {
		return err
	}

	if when, err = when.ApplyTime(speech); err != nil {
		if errors.IsNotATimeError(err) {
			if when.IsItToday() {
				when = calendar.GetNextPeriod(calendar.NewDate(time.Now()))
			} else {
				when = when.Morning()
			}
		}
	}
	reminder.When = when.Until()

	return nil
}

func (obj *DayParser) findDay(speech *speeches.Speech) (*calendar.Date, error) {
	var (
		pattern             *speeches.Pattern
		err                 error
		preposition, moment string
		when                *calendar.Date
		wordsDateMap        = map[dayName]*calendar.Date{
			today:     calendar.NewDate(calendar.LastMidnight()),
			tomorrow:  calendar.NewDate(calendar.NextMidnight()),
			monday:    calendar.NextMonday(),
			tuesday:   calendar.NextTuesday(),
			wednesday: calendar.NextWednesday(),
			thursday:  calendar.NextThursday(),
			friday:    calendar.NextFriday(),
			saturday:  calendar.NextSaturday(),
			sunday:    calendar.NextSunday(),
		}
		ok bool
	)

	if pattern, err = obj.getPattern(speech); err != nil {
		return nil, err
	}

	switch pattern.GetLength() {
	case 2:
		if err = pattern.MakeOut(&preposition, &moment); err != nil {
			return nil, err
		}
	case 1:
		if err = pattern.MakeOut(&moment); err != nil {
			return nil, err
		}
	}

	if when, ok = wordsDateMap[dayName(moment)]; ok == false {
		return nil, errors.NewUnConsistentGlossaryError(pattern.Origin(), obj.Glossary())
	}

	if err = speech.ApplyPattern(pattern); err != nil {
		return nil, err
	}

	return when, nil
}

func (obj *DayParser) getPattern(speech *speeches.Speech) (*speeches.Pattern, error) {
	var (
		pattern             *speeches.Pattern
		err                 error
		preposition, moment string
		prepositions        = map[string]bool{
			`в`:  true,
			`во`: true,
		}
	)

	if pattern, err = speech.TryPattern(2); err == nil {
		if err = pattern.MakeOut(&preposition, &moment); err != nil {
			return nil, err
		}
		if _, ok := prepositions[preposition]; ok {
			return pattern, nil
		}
	}

	return speech.TryPattern(1)
}

func (obj *DayParser) Glossary() []string {
	return []string{
		today.String(),
		tomorrow.String(),
		monday.String(),
		tuesday.String(),
		wednesday.String(),
		thursday.String(),
		friday.String(),
		saturday.String(),
		sunday.String(),
	}
}

func (obj *DayParser) PatternList() []string {
	return []string{
		`утром позавтракать`,
		`завтра позвонить маме`,
		`завтра вечером вынести мусор`,
		`завтра в 15:35 назначить встречу`,
		`в четверг в 16:00 совещание`,
	}
}
