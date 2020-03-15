package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextPeriod(t *testing.T) {
	tests := []struct {
		args *Date
		want *Date
	}{
		{
			args: NewDate(time.Date(2020, 03, 14, 0, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 9, 0, 0, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 0, 1, 1, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 9, 1, 1, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 8, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 9, 0, 0, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 9, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 12, 0, 0, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 12, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 18, 0, 0, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 18, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 14, 23, 0, 0, 0, time.Local)),
		},
		{
			args: NewDate(time.Date(2020, 03, 14, 23, 0, 0, 0, time.Local)),
			want: NewDate(time.Date(2020, 03, 15, 9, 0, 0, 0, time.Local)),
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, GetNextPeriod(tt.args))
	}
}
