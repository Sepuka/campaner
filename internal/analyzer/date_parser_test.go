package analyzer

import (
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDateParser_Parse(t *testing.T) {
	var (
		now                 = time.Now()
		pastDay             = now.Add(-calendar.Day)
		expectedNextYearDay = time.Date(pastDay.Year(), pastDay.Month(), pastDay.Day(), 9, 0, 0, 0, time.Local).Add(calendar.Year)
		tomorrow            = now.Add(calendar.Day)
		tomorrowMorningTime = time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Now().Location()).Add(calendar.Day)
		tomorrowAt1109      = time.Date(now.Year(), now.Month(), now.Day(), 11, 9, 0, 0, time.Local).Add(calendar.Day)
	)

	tests := []struct {
		name             string
		words            []string
		expectedReminder *domain.Reminder
		rest             []string
		wantErr          bool
	}{
		{
			name:             `some past date without a time`,
			words:            []string{pastDay.Format(`01.02`)},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`01.02`), time.Until(expectedNextYearDay)),
			rest:             []string{},
			wantErr:          false,
		},
		{
			name:             `tomorrow without a time`,
			words:            []string{tomorrow.Format(`01.02`)},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`01.02`), time.Until(tomorrowMorningTime)),
			rest:             []string{},
			wantErr:          false,
		},
		{
			name:             `tomorrow with a time`,
			words:            []string{tomorrow.Format(`01.02`), `в`, `11:09`},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`01.02 в 11:09`), time.Until(tomorrowAt1109)),
			rest:             []string{},
			wantErr:          false,
		},
		{
			name:             `tomorrow with incorrect time`,
			words:            []string{tomorrow.Format(`01.02`), `в`, `blah-blah`},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`01.02 в blah-blah`), time.Until(tomorrowMorningTime)),
			rest:             []string{`в`, `blah-blah`},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		actualReminder := &domain.Reminder{}
		obj := NewDateParser()
		actualRest, err := obj.Parse(tt.words, actualReminder)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.rest, actualRest)
		assert.InDelta(t, tt.expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}
