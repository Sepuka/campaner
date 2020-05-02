package speeches

import "github.com/sepuka/campaner/internal/errors"

type Pattern struct {
	value []string
}

func NewPattern(slice []string) *Pattern {
	return &Pattern{
		value: slice,
	}
}

func (p *Pattern) Origin() string {
	return p.value[0]
}

func (p *Pattern) MakeOut(words ...*string) error {
	if len(p.value) != len(words) {
		return errors.NewPatternLengthIncorrect()
	}

	for i, word := range words {
		if word != nil {
			*word = p.value[i]
		}
	}

	return nil
}

func (p *Pattern) GetLength() int {
	return len(p.value)
}
