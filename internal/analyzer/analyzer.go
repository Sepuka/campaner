package analyzer

import (
	"fmt"
	"strings"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"
)

const MaxWordsLength = 100

type Parser interface {
	Parse([]string, *domain.Reminder) ([]string, error)
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
	words := strings.SplitN(text, ` `, MaxWordsLength)
	obj.buildReminder(words, reminder)
}

func (obj *Analyzer) buildReminder(words []string, reminder *domain.Reminder) {
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
		rest = words[1:]
	}

	obj.buildReminder(rest, reminder)
}
