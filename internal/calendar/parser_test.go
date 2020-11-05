package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMorning(t *testing.T) {
	var (
		now      = time.Now()
		expected = NewDate(time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local))
	)
	assert.Equal(t, expected, NewDate(now).Morning())
}

func TestAfternoon(t *testing.T) {
	var (
		now      = time.Now()
		expected = NewDate(time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local))
	)
	assert.Equal(t, expected, NewDate(now).Afternoon())
}

func TestEvening(t *testing.T) {
	var (
		now      = time.Now()
		expected = NewDate(time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.Local))
	)
	assert.Equal(t, expected, NewDate(now).Evening())
}

func TestNight(t *testing.T) {
	var (
		now      = time.Now()
		expected = NewDate(time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.Local))
	)
	assert.Equal(t, expected, NewDate(now).Night())
}
