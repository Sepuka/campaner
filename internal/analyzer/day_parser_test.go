package analyzer

import (
	"reflect"
	"testing"
	"time"

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
		words    []string
		reminder *domain.Reminder
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: `empty sentence`,
			args: args{
				words:    []string{},
				reminder: &domain.Reminder{},
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `any time on tomorrow`,
			args: args{
				words:    []string{`завтра`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowMorningTime)),
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `tomorrow morning`,
			args: args{
				words:    []string{`завтра`, `утром`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowMorningTime)),
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `tomorrow at 11:23 a.m.`,
			args: args{
				words:    []string{`завтра`, `в`, `11:23`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowMorningTime.Add(143*time.Minute))),
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `tomorrow afternoon`,
			args: args{
				words:    []string{`завтра`, `днем`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowAfternoonTime)),
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `tomorrow evening`,
			args: args{
				words:    []string{`завтра`, `вечером`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowEveningTime)),
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: `tomorrow night`,
			args: args{
				words:    []string{`завтра`, `ночью`},
				reminder: domain.NewReminder(0, ``, time.Until(tomorrowNightTime)),
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &DayParser{}
			actualReminder := &domain.Reminder{}
			got, err := obj.Parse(tt.args.words, actualReminder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			assert.InDelta(t, tt.args.reminder.When.Seconds(), actualReminder.When.Seconds(), 1)
		})
	}
}

func TestDayParser_ParseWeekdays(t *testing.T) {
	var (
		mondayMorning = calendar.NextMonday().Add(9 * time.Hour)
	)

	type args struct {
		words    []string
		reminder *domain.Reminder
	}
	tests := []struct {
		name    string
		args    args
		rest    []string
		wantErr bool
	}{
		{
			name: `on Monday`,
			args: args{
				words:    []string{`понедельник`, `встреча`},
				reminder: domain.NewReminder(0, `в понедельник встреча`, mondayMorning.Until()),
			},
			rest:    []string{`встреча`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		obj := &DayParser{}
		actualReminder := &domain.Reminder{}
		got, err := obj.Parse(tt.args.words, actualReminder)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.rest, got)
		assert.InDelta(t, tt.args.reminder.When.Seconds(), actualReminder.When.Seconds(), 1)
	}
}
