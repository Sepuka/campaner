package analyzer

import (
	"github.com/sepuka/campaner/internal/errors"

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
	dayName    string
	DateParser struct {
	}
)

func (dn dayName) String() string {
	return string(dn)
}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (obj *DateParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		offset       = 1
		rest         []string
		word         string
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
	)

	if len(words) == 0 {
		return words, nil
	}
	word = words[0]

	when, ok := wordsDateMap[dayName(word)]
	if ok == false {
		return words, errors.NewUnConsistentGlossaryError(word, obj.Glossary())
	}

	if when, rest, err = when.ApplyTime(words[offset:]); err != nil {
		if errors.GetType(err) == errors.NotATimeError {
			switch dayName(word) {
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

func (obj *DateParser) PatternList() []string {
	return []string{
		`утром позавтракать`,
		`завтра позвонить маме`,
		`завтра вечером вынести мусор`,
		`завтра в 15:35 назначить встречу`,
		`в четверг в 16:00 совещание`,
	}
}
