package analyzer

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/speeches"
)

func TestDateMonthParser(t *testing.T) {
	var (
		now        = time.Now()
		newYearDay = time.Date(now.Year(), 1, 1, 9, 0, 0, 0, time.Local).Add(calendar.Year)
	)

	tests := []struct {
		name     string
		speech   *speeches.Speech
		expected *domain.Reminder
		wantErr  bool
	}{
		{
			name:     `1-th unknown an important event`,
			speech:   speeches.NewSpeech(`1 непонятнонаписаномесяц будет важное событие`),
			expected: &domain.Reminder{},
			wantErr:  true,
		},
		{
			name:   `1-th january an important event`,
			speech: speeches.NewSpeech(`1 января будет важное событие`),
			expected: &domain.Reminder{
				Subject: strings.Split(`будет важное событие`, ` `),
				When:    time.Until(newYearDay),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		actualReminder := &domain.Reminder{}
		obj := NewDateMonthParser()
		err := obj.Parse(tt.speech, actualReminder)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.InDelta(t, tt.expected.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}
