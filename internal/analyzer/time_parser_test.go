package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestOverTimeParser(t *testing.T) {
	parser := NewTimeParser()
	var testCases = map[string]struct {
		words    []string
		rest     []string
		reminder *domain.Reminder
	}{
		`через секунду`: {
			words:    []string{`через`, `секунду`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Second),
		},
		`через одну секунду`: {
			words:    []string{`через`, `1`, `секунду`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Second),
		},
		`через 2 секунды`: {
			words:    []string{`через`, `2`, `секунды`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Second),
		},
		`через 5 сек`: {
			words:    []string{`через`, `5`, `сек`},
			rest:     []string{},
			reminder: domain.NewReminder(0, ``, 5*time.Second),
		},
		`через минуту`: {
			words:    []string{`через`, `минуту`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Minute),
		},
		`через одну минуту`: {
			words:    []string{`через`, `1`, `минуту`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Minute),
		},
		`через 2 минуты`: {
			words:    []string{`через`, `2`, `минуты`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Minute),
		},
		`через час`: {
			words:    []string{`через`, `час`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Hour),
		},
		`через один час`: {
			words:    []string{`через`, `1`, `час`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Hour),
		},
		`через 1.5 часа`: {
			words:    []string{`через`, `1.5`, `часа`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 90*time.Minute),
		},
		`через 2 часа`: {
			words:    []string{`через`, `2`, `часа`, `действие`},
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Hour),
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		rest, err := parser.Parse(testCase.words, actualReminder)
		assert.Equal(t, testCase.rest, rest, testError)
		assert.InDelta(t, testCase.reminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}

func TestOnTimeParser(t *testing.T) {
	parser := NewTimeParser()
	now := time.Now()
	nextDateTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local)
	if now.Hour() > 15 {
		nextDateTime = nextDateTime.Add(24 * time.Hour)
	}

	var testCases = map[string]struct {
		words    []string
		rest     []string
		reminder *domain.Reminder
	}{
		`в 15:00 совершить действие`: {
			words:    []string{`в`, `15:00`, `совершить`, `действие`},
			rest:     []string{`совершить`, `действие`},
			reminder: domain.NewReminder(0, ``, time.Until(nextDateTime)),
		},
		`в 15 часов совершить действие`: {
			words:    []string{`в`, `15`, `часов`, `совершить`, `действие`},
			rest:     []string{`совершить`, `действие`},
			reminder: domain.NewReminder(0, ``, time.Until(nextDateTime)),
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		rest, err := parser.Parse(testCase.words, actualReminder)
		assert.Equal(t, testCase.rest, rest, testError)
		assert.InDelta(t, testCase.reminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}
