package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverTimeParser(t *testing.T) {
	parser := NewTimeParser()
	var testCases = map[string]struct {
		words    []string
		rest     []string
		reminder *Reminder
	}{
		`empty rest when empty words`: {
			words:    []string{},
			rest:     []string{},
			reminder: &Reminder{},
		},
		`stop analyze when occurred unknown word`: {
			words: []string{
				`first_unknown_word`,
				`second_unknown_word`,
			},
			rest: []string{
				`first_unknown_word`,
				`second_unknown_word`,
			},
			reminder: &Reminder{},
		},
		`через секунду`: {
			words: []string{`через`, `секунду`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Second,
			},
		},
		`через одну секунду`: {
			words: []string{`через`, `1`, `секунду`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Second,
			},
		},
		`через 2 секунды`: {
			words: []string{`через`, `2`, `секунды`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: 2 * time.Second,
			},
		},
		`через минуту`: {
			words: []string{`через`, `минуту`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Minute,
			},
		},
		`через одну минуту`: {
			words: []string{`через`, `1`, `минуту`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Minute,
			},
		},
		`через 2 минуты`: {
			words: []string{`через`, `2`, `минуты`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: 2 * time.Minute,
			},
		},
		`через час`: {
			words: []string{`через`, `час`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Hour,
			},
		},
		`через один час`: {
			words: []string{`через`, `1`, `час`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: time.Hour,
			},
		},
		`через 2 часа`: {
			words: []string{`через`, `2`, `часа`, `действие`},
			rest:  []string{`действие`},
			reminder: &Reminder{
				when: 2 * time.Hour,
			},
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &Reminder{}
		rest, err := parser.Parse(testCase.words, actualReminder)
		assert.Equal(t, testCase.rest, rest, testError)
		assert.Equal(t, testCase.reminder, actualReminder)
		assert.NoError(t, err, testError)
	}
}

func TestOnTimeParser(t *testing.T) {
	parser := NewTimeParser()
	now := time.Now()
	day := now.Day()
	if now.Hour() > 15 {
		day++
	}
	nextDateTime := time.Date(now.Year(), now.Month(), day, 15, 0, 0, 0, time.Local)

	var testCases = map[string]struct {
		words    []string
		rest     []string
		reminder *Reminder
	}{
		`в 5 минут совершить действие`: {
			words: []string{`в`, `15:00`, `совершить`, `действие`},
			rest:  []string{`совершить`, `действие`},
			reminder: &Reminder{
				when: time.Since(nextDateTime),
			},
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &Reminder{}
		rest, err := parser.Parse(testCase.words, actualReminder)
		assert.Equal(t, testCase.rest, rest, testError)
		assert.InDelta(t, testCase.reminder.when.Seconds(), actualReminder.when.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}