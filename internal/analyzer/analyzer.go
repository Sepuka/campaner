package analyzer

import (
	"fmt"
	"strings"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"
)

type Parser interface {
	Parse(*speeches.Speech, *domain.Reminder) error
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

func (obj *Analyzer) Analyze(text string, reminder *domain.Reminder) {
	obj.buildReminder(speeches.NewSpeech(text), reminder)
}

func (obj *Analyzer) buildReminder(speech *speeches.Speech, reminder *domain.Reminder) {
	const patternLength = 1
	var (
		err     error
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return
	}

	if parser, ok := obj.glossary[pattern.Origin()]; ok {
		if err = parser.Parse(speech, reminder); err != nil {
			var (
				patterns, what string
			)

			switch errors.GetType(err) {
			case errors.ItIsPastTimeError:
				what = `it is past time!`
			default:
				patterns = strings.Join(parser.PatternList(), "\n")
				what = fmt.Sprintf("use known format, for instance:\n%s\n", patterns)
			}

			*reminder = *domain.NewImmediateReminder(reminder.Whom, what)
			return
		}
	} else {
		if err = speech.ApplyPattern(pattern); err != nil {
			return
		}
	}

	obj.buildReminder(speech, reminder)
}
