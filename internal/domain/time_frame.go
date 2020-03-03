package domain

import (
	"errors"
	"strings"
	"time"
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
	var value, dimension = tf.normalize()

	switch {
	case strings.HasPrefix(dimension, `сек`):
		return time.Duration(value) * time.Second, nil
	case strings.HasPrefix(dimension, `минут`):
		return time.Duration(value) * time.Minute, nil
	case strings.HasPrefix(dimension, `час`):
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, errors.New(`unknown dimension`)
	}
}

func (tf *TimeFrame) normalize() (int, string) {
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
