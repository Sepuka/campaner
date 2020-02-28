package domain

import (
	"testing"
	"time"
)

func TestTimeFrame_GetDuration(t *testing.T) {
	type fields struct {
		value     float64
		dimension string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Duration
		wantErr bool
	}{
		{
			name: `two and half hour`,
			fields: fields{
				value:     2.5,
				dimension: `часа`,
			},
			want:    150 * time.Minute,
			wantErr: false,
		},
		{
			name: `a hour`,
			fields: fields{
				value:     1,
				dimension: `час`,
			},
			want:    time.Hour,
			wantErr: false,
		},
		{
			name: `half an hour`,
			fields: fields{
				value:     0.5,
				dimension: `часа`,
			},
			want:    30 * time.Minute,
			wantErr: false,
		},
		{
			name: `two minutes`,
			fields: fields{
				value:     2,
				dimension: `минуты`,
			},
			want:    2 * time.Minute,
			wantErr: false,
		},
		{
			name: `one and half minute`,
			fields: fields{
				value:     1.5,
				dimension: `минуты`,
			},
			want:    90 * time.Second,
			wantErr: false,
		},
		{
			name: `a minute`,
			fields: fields{
				value:     1,
				dimension: `минута`,
			},
			want:    time.Minute,
			wantErr: false,
		},
		{
			name: `half minute`,
			fields: fields{
				value:     0.5,
				dimension: `минуты`,
			},
			want:    30 * time.Second,
			wantErr: false,
		},
		{
			name: `an one and half second`,
			fields: fields{
				value:     1.5,
				dimension: `секунды`,
			},
			want:    time.Second,
			wantErr: false,
		},
		{
			name: `a second`,
			fields: fields{
				value:     1,
				dimension: `секунда`,
			},
			want:    time.Second,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &TimeFrame{
				value:     tt.fields.value,
				dimension: tt.fields.dimension,
			}
			got, err := tf.GetDuration()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
