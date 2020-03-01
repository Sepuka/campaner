package analyzer

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/sepuka/campaner/internal/domain"
)

var (
	unrecognizedPatterError = errors.New(`cannot recognize pattern`)

	overTime = map[string]bool{
		`через`: true,
	}
	onTime = map[string]bool{
		`в`: true,
	}
	hoursTmpl, _ = regexp.Compile(`\d{1,2}:\d{1,2}`)
)

type TimeParser struct {
	overTime map[string]bool
	onTime   map[string]bool
}

func NewTimeParser() *TimeParser {
	return &TimeParser{
		overTime: overTime,
		onTime:   onTime,
	}
}

func (obj *TimeParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		duration time.Duration
		rest     []string
		err      error
		word     string
	)

	if len(words) == 0 {
		return words, nil
	}

	word = words[0]
	switch {
	case obj.overTime[word]:
		duration, rest, err = obj.overTimeParser(words[1:])
	case obj.onTime[word]:
		duration, rest, err = obj.onTimeParser(words[1:])
	default:
		rest = words
	}

	reminder.When = duration

	return rest, err
}

func (obj *TimeParser) overTimeParser(words []string) (time.Duration, []string, error) {
	var (
		value      float64
		dimension  string
		restOffset int
		err        error
		timeFrame  *domain.TimeFrame
		duration   time.Duration
	)

	if value, err = strconv.ParseFloat(words[0], 32); err != nil {
		value = 1
		dimension = words[0]
		restOffset = 1
	} else {
		dimension = words[1]
		restOffset = 2
	}

	timeFrame = domain.NewTimeFrame(value, dimension)
	if duration, err = timeFrame.GetDuration(); err != nil {
		return duration, words, err
	}

	return duration, words[restOffset:], nil
}

func (obj *TimeParser) onTimeParser(words []string) (time.Duration, []string, error) {
	var (
		exactTime time.Time
		err       error
		current   = time.Now()
		day       = current.Day()
		when      = words[0]
	)

	switch {
	case hoursTmpl.MatchString(when):
		if exactTime, err = time.Parse(`15:04`, when); err != nil {
			return 0, words, unrecognizedPatterError
		}
		if exactTime.Hour() < current.Hour() {
			day++
		}
		exactTime = time.Date(current.Year(), current.Month(), day, exactTime.Hour(), exactTime.Minute(), 0, 0, time.Local)

		return time.Until(exactTime), words[1:], nil
	default:
		return 0, words, unrecognizedPatterError
	}
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
