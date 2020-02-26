package analyzer

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (obj *TimeParser) Parse(words []string, reminder *Reminder) ([]string, error) {
	var (
		duration time.Duration
		rest     []string
		err      error
	)

	if len(words) == 0 {
		return words, nil
	}

	var word = words[0]
	switch {
	case obj.overTime[word]:
		duration, rest, err = obj.overTimeParser(words[1:])
	case obj.onTime[word]:
		duration, rest, err = obj.onTimeParser(words[1:])
	default:
		rest = words
	}

	reminder.when = duration

	return rest, err
}

func (obj *TimeParser) overTimeParser(words []string) (time.Duration, []string, error) {
	var (
		number   int
		quantity string
		idx      int
	)

	number, err := strconv.Atoi(words[0])
	if err != nil {
		number = 1
		quantity = words[0]
		idx = 1
	} else {
		quantity = words[1]
		idx = 2
	}

	switch {
	case strings.HasPrefix(quantity, `секунд`):
		return time.Duration(number) * time.Second, words[idx:], nil
	case strings.HasPrefix(quantity, `минут`):
		return time.Duration(number) * time.Minute, words[idx:], nil
	case strings.HasPrefix(quantity, `час`):
		return time.Duration(number) * time.Hour, words[idx:], nil
	default:
		return 0, words, unrecognizedPatterError
	}
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
