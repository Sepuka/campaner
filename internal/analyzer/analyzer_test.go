package analyzer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/stretchr/testify/assert"
)

var (
	parsers = []Parser{
		NewTimeParser(),
		NewDateParser(),
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
	nextDateTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 15, 0, 0, time.Local)
	if now.After(nextDateTime) {
		nextDateTime = nextDateTime.Add(24 * time.Hour)
	}
	tomorrowMorning := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local).Add(24 * time.Hour)

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
			words:    `напомни В 22:15 что-то сделать`,
			reminder: domain.NewReminder(0, `напомни В 22:15 что-то сделать`, time.Until(nextDateTime)),
		},
		`завтра в 09:23 отвести детей в школу`: {
			words:    `завтра в 09:23 отвести детей в школу`,
			reminder: domain.NewReminder(0, `завтра в 09:23 отвести детей в школу`, time.Until(tomorrowMorning.Add(23*time.Minute))),
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

func TestDayOfWeekAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer(glossary)
	var testCases = map[string]struct {
		words    string
		reminder *domain.Reminder
	}{
		`в понедельник и время с минутами`: {
			words:    `в понедельник в 16:00 часов встреча`,
			reminder: domain.NewReminder(0, `в понедельник в 16:00 часов встреча`, calendar.NextMonday().Add(time.Hour*16).Until()),
		},
		`во вторник и время с минутами`: {
			words:    `во вторник в 17:00 часов встреча`,
			reminder: domain.NewReminder(0, `во вторник в 17:00 часов встреча`, calendar.NextTuesday().Add(time.Hour*17).Until()),
		},
		`в среду и время с минутами`: {
			words:    `в среду в 18:00 часов встреча`,
			reminder: domain.NewReminder(0, `в среду в 18:00 часов встреча`, calendar.NextWednesday().Add(time.Hour*18).Until()),
		},
		`в четверг и время с минутами`: {
			words:    `в четверг в 19:00 часов встреча`,
			reminder: domain.NewReminder(0, `в четверг в 19:00 часов встреча`, calendar.NextThursday().Add(time.Hour*19).Until()),
		},
		`в пятницу и время с минутами`: {
			words:    `в пятницу в 20:00 часов встреча`,
			reminder: domain.NewReminder(0, `в пятницу в 20:00 часов встреча`, calendar.NextFriday().Add(time.Hour*20).Until()),
		},
		`в субботу и время с минутами`: {
			words:    `в субботу в 21:00 час встреча`,
			reminder: domain.NewReminder(0, `в субботу в 21:00 час встреча`, calendar.NextSaturday().Add(time.Hour*21).Until()),
		},
		`в воскресенье и время с минутами`: {
			words:    `в воскресенье в 22:00 часа встреча`,
			reminder: domain.NewReminder(0, `в воскресенье в 22:00 часа встреча`, calendar.NextSunday().Add(time.Hour*22).Until()),
		},
		`день недели и время без минут`: {
			words:    `В среду в 16 пройдет маленькая пятничная встреча.`,
			reminder: domain.NewReminder(0, `В среду в 16 пройдет маленькая пятничная встреча.`, calendar.NextWednesday().Add(time.Hour*16).Until()),
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
