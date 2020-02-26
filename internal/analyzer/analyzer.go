package analyzer

import (
	"fmt"
	"strings"
)

const MaxWordsLength = 100

type Parser interface {
	Parse([]string, *Reminder) ([]string, error)
	Glossary() []string
	PatternList() []string
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

func (obj *Analyzer) Analyze(text string, reminder *Reminder) {
	words := strings.SplitN(text, ` `, MaxWordsLength)
	obj.buildReminder(words, reminder)

	if reminder.what == `` {
		reminder.what = text
	}
}

func (obj *Analyzer) buildReminder(words []string, reminder *Reminder) {
	if len(words) == 0 {
		return
	}

	var (
		rest []string
		err  error
	)
	if parser, ok := obj.glossary[words[0]]; ok {
		if rest, err = parser.Parse(words, reminder); err != nil {
			var (
				patterns = strings.Join(parser.PatternList(), "\n")
				what     = fmt.Sprintf("use known format, for instance: %s\n", patterns)
			)
			reminder = NewImmediateReminder(reminder.Whom(), what)
			return
		}
	} else {
		rest = words[1:]
	}

	obj.buildReminder(rest, reminder)
}
