package domain

import (
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/errors"
)

const (
	maxMinuteValue = 59
	maxHourValue   = 23
)

type TimeFrame struct {
	value     float64
	dimension string
}

func NewTimeFrame(value float64, dimension string) *TimeFrame {
	return &TimeFrame{
		value:     value,
		dimension: dimension,
	}
}

func (tf *TimeFrame) GetDuration() (time.Duration, error) {
	var value, dimension = tf.castToWhole()

	switch {
	case strings.HasPrefix(dimension, `сек`):
		return time.Duration(value) * time.Second, nil
	case strings.HasPrefix(dimension, `минут`):
		return time.Duration(value) * time.Minute, nil
	case strings.HasPrefix(dimension, `час`):
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, errors.NewUnknownDimensionError(tf.value, tf.dimension)
	}
}

func (tf *TimeFrame) GetTime() (time.Time, error) {
	var (
		now    = time.Now()
		value  = int(tf.value)
		atTime time.Time
	)

	switch {
	case strings.HasPrefix(tf.dimension, `минут`):
		if value > maxMinuteValue {
			return time.Time{}, errors.NewInvalidTimeValueError(value)
		}
		atTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), value, 0, 0, time.Local)
		if atTime.Before(time.Now()) {
			atTime = atTime.Add(time.Hour)
		}
		return atTime, nil
	case strings.HasPrefix(tf.dimension, `час`):
		if value > maxHourValue {
			return time.Time{}, errors.NewInvalidTimeValueError(value)
		}
		atTime = time.Date(now.Year(), now.Month(), now.Day(), value, 0, 0, 0, time.Local)
		if atTime.Before(time.Now()) {
			atTime = atTime.Add(24 * time.Hour)
		}
		return atTime, nil
	default:
		return time.Time{}, errors.NewUnknownDimensionError(tf.value, tf.dimension)
	}
}

func (tf *TimeFrame) castToWhole() (int, string) {
	var (
		value     int
		dimension = tf.dimension
	)

	if tf.isValueInt() {
		return int(tf.value), tf.dimension
	}

	switch {
	case strings.HasPrefix(tf.dimension, `час`):
		value = int(tf.value * 60)
		dimension = `минут`
	case strings.HasPrefix(tf.dimension, `минут`):
		value = int(tf.value * 60)
		dimension = `секунд`
	default:
		value = int(tf.value)
	}

	return value, dimension
}

func (tf *TimeFrame) isValueInt() bool {
	return float64(int64(tf.value)) == tf.value
}
