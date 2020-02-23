package analyzer

import (
	"fmt"
	"os"
	"testing"
	"time"

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
		reminder *Reminder
	}{
		`empty rest when empty words`: {
			words:    ``,
			reminder: &Reminder{},
		},
		`напомни мне через 25 секунд что-то сделать`: {
			words:    `напомни мне через 25 секунд что-то сделать`,
			reminder: &Reminder{when: time.Duration(25) * time.Second},
		},
		`напомни в 22:15 что-то сделать`: {
			words:    `напомни в 22:15 что-то сделать`,
			reminder: &Reminder{when: time.Since(nextDateTime)},
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder, err := analyzer.Analyze(testCase.words)
		assert.InDelta(t, testCase.reminder.when.Seconds(), actualReminder.when.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}