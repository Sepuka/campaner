package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDateParser(t *testing.T) {
	var (
		now = time.Now()

		pastDay         = now.Add(-calendar.Day)
		expectedPastDay = time.Date(pastDay.Year(), pastDay.Month(), pastDay.Day(), 9, 0, 0, 0, time.Local).Add(calendar.Year)

		expectedToday = calendar.GetNextPeriod(calendar.NewDate(now))

		tomorrow            = now.Add(calendar.Day)
		tomorrowMorningTime = time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Now().Location()).Add(calendar.Day)
		tomorrowAt1109      = time.Date(now.Year(), now.Month(), now.Day(), 11, 9, 0, 0, time.Local).Add(calendar.Day)
	)

	tests := []struct {
		name     string
		speech   *speeches.Speech
		expected *domain.Reminder
		wantErr  bool
	}{
		{
			name:   `some alone past date without a time`,
			speech: speeches.NewSpeech(pastDay.Format(calendar.DayMonthFormat)),
			expected: &domain.Reminder{
				Subject: []string{pastDay.Format(calendar.DayMonthFormat)},
				When:    time.Until(expectedPastDay),
			},
			wantErr: false,
		},
		{
			name:   `some alone current date without a time`,
			speech: speeches.NewSpeech(now.Format(calendar.DayMonthFormat)),
			expected: &domain.Reminder{
				Subject: []string{now.Format(calendar.DayMonthFormat)},
				When:    expectedToday.Until(),
			},
			wantErr: false,
		},
		{
			name:   `tomorrow without a time`,
			speech: speeches.NewSpeech(tomorrow.Format(calendar.DayMonthFormat)),
			expected: &domain.Reminder{
				Subject: []string{tomorrow.Format(calendar.DayMonthFormat)},
				When:    time.Until(tomorrowMorningTime),
			},
			wantErr: false,
		},
		{
			name:   `tomorrow with a time`,
			speech: speeches.NewSpeech(fmt.Sprintf(`%s в 11:09`, tomorrow.Format(calendar.DayMonthFormat))),
			expected: &domain.Reminder{
				Subject: []string{tomorrow.Format(`02.01 в 11:09`)},
				When:    time.Until(tomorrowAt1109),
			},
			wantErr: false,
		},
		{
			name:   `tomorrow with incorrect time`,
			speech: speeches.NewSpeech(fmt.Sprintf(`%s в blah-blah`, tomorrow.Format(calendar.DayMonthFormat))),
			expected: &domain.Reminder{
				Subject: []string{tomorrow.Format(`02.01 в blah-blah`)},
				When:    time.Until(tomorrowMorningTime),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		actualReminder := &domain.Reminder{}
		obj := NewDateParser()
		err := obj.Parse(tt.speech, actualReminder)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.InDelta(t, tt.expected.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}
