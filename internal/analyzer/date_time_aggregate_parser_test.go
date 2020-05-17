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

	parser := NewDateTimeAggregateParser([]Parser{
		NewTimeParser(),
		NewDayParser(),
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
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actualReminder := &domain.Reminder{}
		err := parser.Parse(testCase.speech, actualReminder)
		assert.InDelta(t, testCase.expected.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.NoError(t, err, testError)
	}
}
