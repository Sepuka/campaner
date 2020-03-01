package analyzer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/stretchr/testify/assert"
)

var (
	parsers = []Parser{
		NewTimeParser(),
	}
	glossary = make(map[string]Parser)
)

func buildGlossary() {
	var (
		keyword string
	)
	for _, parser := range parsers {
		for _, keyword = range parser.(Parser).Glossary() {
			glossary[keyword] = parser.(Parser)
		}
	}
}
func TestMain(m *testing.M) {
	buildGlossary()
	code := m.Run()
	os.Exit(code)
}

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer(glossary)
	now := time.Now()
	day := now.Day()
	if now.Hour() > 22 && now.Minute() > 15 {
		day++
	}
	nextDateTime := time.Date(now.Year(), now.Month(), day, 22, 15, 0, 0, time.Local)

	var testCases = map[string]struct {
		words    string
		reminder *domain.Reminder
	}{
		`empty rest when empty words`: {
			words:    ``,
			reminder: &domain.Reminder{},
		},
		`unknown pattern`: {
			words:    `abc`,
			reminder: domain.NewReminder(0, `abc`, time.Nanosecond),
		},
		`напомни мне через 25 секунд что-то сделать`: {
			words:    `напомни мне через 25 секунд что-то сделать`,
			reminder: domain.NewReminder(0, `напомни мне через 25 секунд что-то сделать`, time.Duration(25)*time.Second),
		},
		`напомни в 22:15 что-то сделать`: {
			words:    `напомни в 22:15 что-то сделать`,
			reminder: domain.NewReminder(0, `напомни в 22:15 что-то сделать`, time.Until(nextDateTime)),
		},
	}

	for testName, testCase := range testCases {
		var (
			testError        = fmt.Sprintf(`test "%s" error`, testName)
			expectedReminder = testCase.reminder
			actualReminder   = domain.NewReminder(0, testCase.words, time.Nanosecond)
		)
		analyzer.Analyze(testCase.words, actualReminder)
		assert.InDelta(t, expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.Equal(t, expectedReminder.What, actualReminder.What, testError)
	}
}
