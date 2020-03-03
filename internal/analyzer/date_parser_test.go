package analyzer

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sepuka/campaner/internal/domain"
)

func TestDateParser_Parse(t *testing.T) {
	now := time.Now()
	day := now.Day()
	if now.Hour() > 15 {
		day++
	}
	tomorrowMorningTime := time.Date(now.Year(), now.Month(), day, 9, 0, 0, 0, time.Now().Location())
	tomorrowAfternoonTime := time.Date(now.Year(), now.Month(), day, 12, 0, 0, 0, time.Now().Location())
	tomorrowEveningTime := time.Date(now.Year(), now.Month(), day, 18, 0, 0, 0, time.Now().Location())
	tomorrowNightTime := time.Date(now.Year(), now.Month(), day, 23, 0, 0, 0, time.Now().Location())

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
			obj := &DateParser{}
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
