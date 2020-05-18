package calendar

import (
	"strconv"
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/errors"
)

func findTimeOfADay(date *Date, speech *speeches.Speech) (*Date, error) {
	var (
		pattern *speeches.Pattern
		err     error
	)

	if pattern, err = speech.TryPattern(1); err != nil {
		return date, err
	}

	switch pattern.Origin() {
	case `утром`:
		date = date.Morning()
	case `днем`:
		date = date.Afternoon()
	case `вечером`:
		date = date.Evening()
	case `ночью`:
		date = date.Night()
	default:
		return date, errors.NewNotATimeError()
	}

	return date, speech.ApplyPattern(pattern)
}

func onTimeParser(date *Date, speech *speeches.Speech) (*Date, error) {
	const patternLength = 2
	var (
		atTime              time.Time
		err                 error
		preposition, moment string
		pattern             *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return date, err
	}

	if err = pattern.MakeOut(&preposition, &moment); err != nil {
		return date, err
	}

	if preposition != `в` {
		return date, errors.NewNotATimeError()
	}

	if atTime, err = time.Parse(HourMinuteSecondFormat, moment); err != nil {
		if atTime, err = time.Parse(HourMinuteFormat, moment); err != nil {
			return date, errors.NewNotATimeError()
		}
	}

	atTime = time.Date(date.date.Year(), date.date.Month(), date.date.Day(), atTime.Hour(), atTime.Minute(), 0, 0, time.Local)

	if atTime.Before(time.Now()) {
		atTime = atTime.Add(Day)
	}

	if err = speech.ApplyPattern(pattern); err != nil {
		return nil, err
	}

	if pattern, err = speech.TryPattern(1); err != nil {
		if errors.IsSpeechIsOverError(err) {
			return NewDate(atTime), nil
		}
		return nil, err
	}

	if strings.HasPrefix(pattern.Origin(), `час`) {
		return NewDate(atTime), speech.ApplyPattern(pattern)
	}

	return NewDate(atTime), nil
}

func overTimeParser(date *Date, speech *speeches.Speech) (*Date, error) {
	const (
		shortPatternLength = 2
		fullPatternLength  = 3
	)
	var (
		value                          float64
		err                            error
		timeFrame                      *domain.TimeFrame
		duration                       time.Duration
		preposition, moment, dimension string
		pattern                        *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(fullPatternLength); err != nil {
		if pattern, err = speech.TryPattern(shortPatternLength); err != nil {
			return date, err
		} else {
			if err = pattern.MakeOut(&preposition, &dimension); err != nil {
				return date, err
			} else {
				moment = `1`
			}
		}
	} else {
		if err = pattern.MakeOut(&preposition, &moment, &dimension); err != nil {
			return date, err
		}
	}

	if preposition != `через` {
		return nil, errors.NewNotATimeError()
	}

	if value, err = strconv.ParseFloat(moment, 32); err != nil {
		value = 1
		dimension = moment
	}

	timeFrame = domain.NewTimeFrame(value, dimension)
	if duration, err = timeFrame.GetDuration(); err != nil {
		return nil, errors.NewNotATimeError()
	}

	return NewDate(time.Now().Add(duration)), speech.ApplyPattern(pattern)
}

func onNextTimePeriod(date *Date, speech *speeches.Speech) (*Date, error) {
	const patternLength = 3
	var (
		atTime            time.Time
		err               error
		moment, dimension string
		value             float64
		timeFrame         *domain.TimeFrame
		pattern           *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return date, err
	}

	if pattern.Origin() != `в` {
		return nil, errors.NewNotATimeError()
	}

	if err = pattern.MakeOut(nil, &moment, &dimension); err != nil {
		return date, err
	}

	if value, err = strconv.ParseFloat(moment, 9); err != nil {
		return nil, errors.NewNotATimeError()
	}

	timeFrame = domain.NewTimeFrame(value, dimension)
	if atTime, err = timeFrame.GetTime(); err != nil {
		return nil, errors.NewNotATimeError()
	}

	return NewDate(atTime), speech.ApplyPattern(pattern)
}

func onNextTimeHourPeriod(date *Date, speech *speeches.Speech) (*Date, error) {
	const patternLength = 2
	var (
		atTime              *Date
		err                 error
		value               int64
		pattern             *speeches.Pattern
		preposition, moment string
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return date, err
	}

	if err = pattern.MakeOut(&preposition, &moment); err != nil {
		return date, err
	}

	if preposition != `в` {
		return nil, errors.NewNotATimeError()
	}

	if value, err = strconv.ParseInt(moment, 0, 8); err != nil {
		return date, errors.NewNotATimeError()
	}

	if value < 0 || value > 23 {
		return date, errors.NewNotATimeError()
	}

	atTime = date.Add(time.Hour * time.Duration(value))
	if atTime.IsPast() {
		atTime = atTime.Add(Day)
	}

	return atTime, speech.ApplyPattern(pattern)
}