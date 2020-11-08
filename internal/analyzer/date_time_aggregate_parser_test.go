package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/speeches"
	"github.com/stretchr/testify/assert"
)

func TestDateTimeAggregateParser(t *testing.T) {
	expectedTime := calendar.NewDate(calendar.LastMidnight()).Add(16 * time.Hour)
	if expectedTime.IsPast() {
		expectedTime = expectedTime.Add(calendar.Day)
	}
	expectedOver2Days := calendar.NewDate(time.Now()).Morning().Add(2 * calendar.Day)
	expectedOver5Days := calendar.NewDate(time.Now()).Morning().Add(5 * calendar.Day)
	expectedOver2Months := calendar.NewDate(calendar.LastMidnight().Add(9 * time.Hour)).Add(calendar.Day * 60)
	expectedOver8Weeks := calendar.NewDate(calendar.LastMidnight()).Add(9 * time.Hour).Add(calendar.Day * 7 * 8)

	parser := NewDateTimeAggregateParser([]Parser{
		NewTimeParser(),
		NewDayParser(),
		NewPeriodParser(),
	})
	var testCases = map[string]struct {
		speech   *speeches.Speech
		expected *domain.Reminder
	}{
		`в среду`: {
			speech: speeches.NewSpeech(`в среду будет вечеринка`),
			expected: &domain.Reminder{
				When: calendar.NextWednesday().Morning().Until(),
			},
		},
		`в 16:00`: {
			speech: speeches.NewSpeech(`в 16:00 наступит встреча`),
			expected: &domain.Reminder{
				When: expectedTime.Until(),
			},
		},
		`в школу`: {
			speech: speeches.NewSpeech(`отвезти детей в школу`),
			expected: &domain.Reminder{
				When: time.Microsecond,
			},
		},
		`через 2 дня`: {
			speech: speeches.NewSpeech(`через 2 дня напомни`),
			expected: &domain.Reminder{
				When: expectedOver2Days.Until(),
			},
		},
		`через 5 дней`: {
			speech: speeches.NewSpeech(`через 5 дней будет праздник`),
			expected: &domain.Reminder{
				When: expectedOver5Days.Until(),
			},
		},
		`через 8 недель`: {
			speech: speeches.NewSpeech(`через 8 недель будет что-то`),
			expected: &domain.Reminder{
				When: expectedOver8Weeks.Until(),
			},
		},
		`через 2 месяца`: {
			speech: speeches.NewSpeech(`через 2 месяца новый год`),
			expected: &domain.Reminder{
				When: expectedOver2Months.Until(),
			},
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		err := parser.Parse(testCase.speech, actualReminder)
		assert.InDelta(t, testCase.expected.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err)
	}
}
