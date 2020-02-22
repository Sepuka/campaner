package analyzer

import (
	"strings"
)

const MaxWordsLength = 100

type Parser interface {
	Parse([]string, *Reminder) ([]string, error)
	Glossary() []string
}

type Glossary map[string]Parser

type Analyzer struct {
	glossary Glossary
}

func NewAnalyzer(glossary Glossary) *Analyzer {
	return &Analyzer{
		glossary: glossary,
	}
}

func (obj *Analyzer) Analyze(text string, reminder *Reminder) error {
	var (
		err error
	)

	words := strings.SplitN(text, ` `, MaxWordsLength)
	err = obj.buildReminder(words, reminder)
	reminder.what = text

	return err
}

func (obj *Analyzer) buildReminder(words []string, reminder *Reminder) error {
	if len(words) == 0 {
		return nil
	}

	var (
		rest []string
		err  error
	)
	if parser, ok := obj.glossary[words[0]]; ok {
		if rest, err = parser.Parse(words, reminder); err != nil {
			return err
		}
	} else {
		rest = words[1:]
	}

	return obj.buildReminder(rest, reminder)
}
