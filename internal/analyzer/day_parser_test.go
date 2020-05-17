package analyzer

import (
	"strings"
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/stretchr/testify/assert"

	"github.com/sepuka/campaner/internal/domain"
)

func TestDayParser_Parse(t *testing.T) {
	now := time.Now()
	tomorrowMorningTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Now().Location()).Add(24 * time.Hour)
	tomorrowAfternoonTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Now().Location()).Add(24 * time.Hour)
	tomorrowEveningTime := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.Now().Location()).Add(24 * time.Hour)
	tomorrowNightTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.Now().Location()).Add(24 * time.Hour)

	type args struct {
		speech   *speeches.Speech
		reminder *domain.Reminder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: `any time on tomorrow`,
			args: args{
				speech: speeches.NewSpeech(`завтра`),
				reminder: &domain.Reminder{
					Subject: []string{`завтра`},
					When:    time.Until(tomorrowMorningTime),
				},
			},
		},
		{
			name: `tomorrow morning`,
			args: args{
				speech: speeches.NewSpeech(`завтра утром`),
				reminder: &domain.Reminder{
					Subject: []string{`завтра`},
					When:    time.Until(tomorrowMorningTime),
				},
			},
		},
		{
			name: `tomorrow at 11:23 a.m.`,
			args: args{
				speech: speeches.NewSpeech(`завтра в 11:23`),
				reminder: &domain.Reminder{
					Subject: strings.Split(`завтра в 11:23`, ` `),
					When:    time.Until(tomorrowMorningTime.Add(143 * time.Minute)),
				},
			},
		},
		{
			name: `tomorrow afternoon`,
			args: args{
				speech: speeches.NewSpeech(`завтра днем`),
				reminder: &domain.Reminder{
					Subject: strings.Split(`завтра днем`, ` `),
					When:    time.Until(tomorrowAfternoonTime),
				},
			},
		},
		{
			name: `tomorrow evening`,
			args: args{
				speech: speeches.NewSpeech(`завтра вечером`),
				reminder: &domain.Reminder{
					Subject: strings.Split(`завтра вечером`, ` `),
					When:    time.Until(tomorrowEveningTime),
				},
			},
		},
		{
			name: `tomorrow night`,
			args: args{
				speech: speeches.NewSpeech(`завтра ночью`),
				reminder: &domain.Reminder{
					Subject: strings.Split(`завтра ночью`, ` `),
					When:    time.Until(tomorrowNightTime),
				},
			},
		},
	}
	for _, tt := range tests {
		obj := NewDayParser()
		actualReminder := &domain.Reminder{}
		err := obj.Parse(tt.args.speech, actualReminder)
		assert.NoError(t, err)
		assert.InDelta(t, tt.args.reminder.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}

func TestDayParser_ParseWeekdays(t *testing.T) {
	var (
		mondayMorning = calendar.NextMonday().Add(9 * time.Hour)
	)

	type args struct {
		speech   *speeches.Speech
		reminder *domain.Reminder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: `on Monday`,
			args: args{
				speech: speeches.NewSpeech(`понедельник встреча`),
				reminder: &domain.Reminder{
					Subject: strings.Split(`в понедельник встреча`, ` `),
					When:    mondayMorning.Until(),
				},
			},
		},
	}

	for _, tt := range tests {
		obj := NewDayParser()
		actualReminder := &domain.Reminder{}
		err := obj.Parse(tt.args.speech, actualReminder)
		assert.NoError(t, err)
		assert.InDelta(t, tt.args.reminder.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}
