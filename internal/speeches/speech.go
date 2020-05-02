package speeches

import (
	"strings"
	"sync"

	"github.com/sepuka/campaner/internal/errors"
)

const (
	separator = ` `
)

type Speech struct {
	sync.Mutex
	pointer  int
	words    []string
	original string
}

func NewSpeech(text string) *Speech {
	if len(text) == 0 {
		return &Speech{
			words: []string{},
		}
	}

	words := strings.Split(strings.ToLower(text), separator)

	return &Speech{
		words:    words,
		original: text,
	}
}

func (s *Speech) IsTheEnd() bool {
	return len(s.words) == 0 || len(s.words)-1 == s.pointer
}

func (s *Speech) TryPattern(length int) (*Pattern, error) {
	if len(s.words) < s.pointer+length {
		return nil, errors.NewSpeechIsOverError(s.pointer, s.original)
	}

	return NewPattern(s.words[s.pointer : s.pointer+length]), nil
}

func (s *Speech) ApplyPattern(pattern *Pattern) error {
	var (
		length         = pattern.GetLength()
		wantPointerPos = s.pointer + length
	)

	s.Lock()
	defer s.Unlock()

	if len(s.words) < wantPointerPos {
		return errors.NewSpeechIsOverError(s.pointer, s.original)
	}

	s.pointer += length

	return nil
}
