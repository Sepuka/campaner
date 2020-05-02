package speeches

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpeech(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want *Speech
	}{
		{
			name: `empty speech`,
			args: args{
				text: ``,
			},
			want: &Speech{
				words:    []string{},
				pointer:  0,
				original: ``,
			},
		},
		{
			name: `one word`,
			args: args{
				text: `word`,
			},
			want: &Speech{
				words:    []string{`word`},
				pointer:  0,
				original: `word`,
			},
		},
		{
			name: `two words`,
			args: args{
				text: `two words`,
			},
			want: &Speech{
				words:    []string{`two`, `words`},
				pointer:  0,
				original: `two words`,
			},
		},
	}
	for _, tt := range tests {
		got := NewSpeech(tt.args.text)
		assert.Equal(t, tt.want.original, got.original)
		assert.Equal(t, tt.want.words, got.words)
		assert.Equal(t, 0, got.pointer)
	}
}

func TestSpeech_GetPattern(t *testing.T) {
	tests := []struct {
		name          string
		speech        *Speech
		patternLength int
		want          *Pattern
		wantErr       bool
	}{
		{
			name:          `empty speech`,
			speech:        NewSpeech(``),
			patternLength: 1,
			want:          nil,
			wantErr:       true,
		},
		{
			name:          `one-word speech with pattern length is 1`,
			speech:        NewSpeech(`word`),
			want:          NewPattern([]string{`word`}),
			patternLength: 1,
			wantErr:       false,
		},
		{
			name:          `one-word speech with pattern length is 2`,
			speech:        NewSpeech(`word`),
			want:          nil,
			patternLength: 2,
			wantErr:       true,
		},
		{
			name:          `two-words speech with pattern length is 1`,
			speech:        NewSpeech(`two words`),
			patternLength: 1,
			want:          NewPattern([]string{`two`}),
			wantErr:       false,
		},
		{
			name:          `two-words speech with pattern length is 2`,
			speech:        NewSpeech(`two words`),
			patternLength: 2,
			want:          NewPattern([]string{`two`, `words`}),
			wantErr:       false,
		},
		{
			name:          `two-words speech with pattern length is 3`,
			speech:        NewSpeech(`two words`),
			patternLength: 3,
			want:          nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		got, err := tt.speech.TryPattern(tt.patternLength)
		if tt.wantErr == true {
			assert.True(t, err != nil)
		} else {
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestSpeech_IsTheEnd(t *testing.T) {
	tests := []struct {
		name   string
		speech *Speech
		want   bool
	}{
		{
			name:   `empty speech`,
			speech: NewSpeech(``),
			want:   true,
		},
		{
			name:   `one word`,
			speech: NewSpeech(`word`),
			want:   true,
		},
		{
			name:   `two words`,
			speech: NewSpeech(`two words`),
			want:   false,
		},
		{
			name: `two words`,
			speech: &Speech{
				words:   []string{`two`, `words`},
				pointer: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		got := tt.speech.IsTheEnd()
		assert.Equal(t, tt.want, got)
	}
}

func TestSpeech_ApplyPattern(t *testing.T) {
	tests := []struct {
		speech          *Speech
		pattern         *Pattern
		expectedPointer int
		wantErr         bool
	}{
		{
			speech:          NewSpeech(``),
			pattern:         NewPattern([]string{}),
			expectedPointer: 0,
			wantErr:         false,
		},
		{
			speech:          NewSpeech(`word`),
			pattern:         NewPattern([]string{`word`}),
			expectedPointer: 1,
			wantErr:         false,
		},
		{
			speech:          NewSpeech(`two word`),
			pattern:         NewPattern([]string{`two`, `word`}),
			expectedPointer: 2,
			wantErr:         false,
		},
		{
			speech:          NewSpeech(`word`),
			pattern:         NewPattern([]string{`two`, `word`}),
			expectedPointer: 0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		err := tt.speech.ApplyPattern(tt.pattern)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.expectedPointer, tt.speech.pointer)
	}
}
