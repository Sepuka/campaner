package analyzer

import (
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDateParser_Past(t *testing.T) {
	var (
		now         = time.Now()
		pastDay     = now.Add(-calendar.Day)
		expectedDay = time.Date(pastDay.Year(), pastDay.Month(), pastDay.Day(), 9, 0, 0, 0, time.Local).Add(calendar.Year)
	)

	tests := []struct {
		name     string
		words    []string
		expected *domain.Reminder
		rest     []string
		wantErr  bool
	}{
		{
			name:     `some alone past date without a time`,
			words:    []string{pastDay.Format(calendar.DayMonthFormat)},
			expected: domain.NewReminder(0, pastDay.Format(calendar.DayMonthFormat), time.Until(expectedDay)),
			rest:     []string{},
			wantErr:  false,
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
		assert.InDelta(t, tt.expected.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}

func TestDateParser_Today(t *testing.T) {
	var (
		now         = time.Now()
		expectedDay = calendar.GetNextPeriod(calendar.NewDate(now))
	)

	tests := []struct {
		name     string
		words    []string
		expected *domain.Reminder
		rest     []string
		wantErr  bool
	}{
		{
			name:     `some alone today date without a time`,
			words:    []string{now.Format(calendar.DayMonthFormat)},
			expected: domain.NewReminder(0, now.Format(calendar.DayMonthFormat), expectedDay.Until()),
			rest:     []string{},
			wantErr:  false,
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
		assert.InDelta(t, tt.expected.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}

func TestDateParser_Future(t *testing.T) {
	var (
		now                 = time.Now()
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
			name:             `tomorrow without a time`,
			words:            []string{tomorrow.Format(calendar.DayMonthFormat)},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(calendar.DayMonthFormat), time.Until(tomorrowMorningTime)),
			rest:             []string{},
			wantErr:          false,
		},
		{
			name:             `tomorrow with a time`,
			words:            []string{tomorrow.Format(calendar.DayMonthFormat), `в`, `11:09`},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`02.01 в 11:09`), time.Until(tomorrowAt1109)),
			rest:             []string{},
			wantErr:          false,
		},
		{
			name:             `tomorrow with incorrect time`,
			words:            []string{tomorrow.Format(calendar.DayMonthFormat), `в`, `blah-blah`},
			expectedReminder: domain.NewReminder(0, tomorrow.Format(`02.01 в blah-blah`), time.Until(tomorrowMorningTime)),
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
