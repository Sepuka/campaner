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
	const patternLength = 1
	var (
		pattern      *speeches.Pattern
		err          error
		when         *calendar.Date
		wordsDateMap = map[dayName]*calendar.Date{
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

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	if when, ok = wordsDateMap[dayName(pattern.Origin())]; ok == false {
		return errors.NewUnConsistentGlossaryError(pattern.Origin(), obj.Glossary())
	}

	if err = speech.ApplyPattern(pattern); err != nil {
		return err
	}

	if when, err = when.ApplyTime(speech); err != nil {
		if errors.IsNotATimeError(err) {
			switch dayName(pattern.Origin()) {
			case today:
				when = calendar.GetNextPeriod(calendar.NewDate(time.Now()))
			default:
				when = when.Morning()
			}
		}
	}
	reminder.When = when.Until()

	return nil
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
