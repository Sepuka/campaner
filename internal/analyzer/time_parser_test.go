package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestOverTimeParser(t *testing.T) {
	parser := NewTimeParser()
	var testCases = map[string]struct {
		speech   *speeches.Speech
		rest     []string
		reminder *domain.Reminder
	}{
		`через секунду`: {
			speech:   speeches.NewSpeech(`через секунду действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Second),
		},
		`через одну секунду`: {
			speech:   speeches.NewSpeech(`через 1 секунду действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Second),
		},
		`через 2 секунды`: {
			speech:   speeches.NewSpeech(`через 2 секунды действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Second),
		},
		`через 5 сек`: {
			speech:   speeches.NewSpeech(`через 5 сек`),
			rest:     []string{},
			reminder: domain.NewReminder(0, ``, 5*time.Second),
		},
		`через минуту`: {
			speech:   speeches.NewSpeech(`через минуту действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Minute),
		},
		`через одну минуту`: {
			speech:   speeches.NewSpeech(`через 1 минуту действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Minute),
		},
		`через 2 минуты`: {
			speech:   speeches.NewSpeech(`через 2 минуты действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Minute),
		},
		`через час`: {
			speech:   speeches.NewSpeech(`через час действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Hour),
		},
		`через один час`: {
			speech:   speeches.NewSpeech(`через 1 час действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, time.Hour),
		},
		`через 1.5 часа`: {
			speech:   speeches.NewSpeech(`через 1.5 часа действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 90*time.Minute),
		},
		`через 2 часа`: {
			speech:   speeches.NewSpeech(`через 2 часа действие`),
			rest:     []string{`действие`},
			reminder: domain.NewReminder(0, ``, 2*time.Hour),
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		err := parser.Parse(testCase.speech, actualReminder)
		assert.InDelta(t, testCase.reminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}

func TestOnTimeParser(t *testing.T) {
	parser := NewTimeParser()
	now := time.Now()
	nextDateTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local)
	if now.Hour() >= 15 {
		nextDateTime = nextDateTime.Add(24 * time.Hour)
	}

	var testCases = map[string]struct {
		speech   *speeches.Speech
		rest     []string
		reminder *domain.Reminder
	}{
		`в 15:00 совершить действие`: {
			speech:   speeches.NewSpeech(`в 15:00 совершить действие`),
			rest:     []string{`совершить`, `действие`},
			reminder: domain.NewReminder(0, ``, time.Until(nextDateTime)),
		},
		`в 15 часов совершить действие`: {
			speech:   speeches.NewSpeech(`в 15 часов совершить действие`),
			rest:     []string{`совершить`, `действие`},
			reminder: domain.NewReminder(0, ``, time.Until(nextDateTime)),
		},
		`в 15 совещание`: {
			speech:   speeches.NewSpeech(`в 15 совещание`),
			reminder: domain.NewReminder(0, ``, time.Until(nextDateTime)),
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		err := parser.Parse(testCase.speech, actualReminder)
		assert.InDelta(t, testCase.reminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}
