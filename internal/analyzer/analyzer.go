package analyzer

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

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

func (a *Analyzer) Analyze(text string, reminder *domain.Reminder) {
	a.buildReminder(speeches.NewSpeech(text), reminder)
}

func (a *Analyzer) buildReminder(speech *speeches.Speech, reminder *domain.Reminder) {
	const patternLength = 1
	var (
		err     error
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		if reminder.GetSubject() == `` {
			reminder.AppendSubject(speeches.NewPattern([]string{`ring!`}))
		}
		if reminder.IsTimeUnknown() {
			var randomSubject = fmt.Sprintf(`Попробуйте фразу: "%s"`, a.getRandomStatement(int64(reminder.Whom)))
			reminder.RewriteSubject(randomSubject)
			reminder.When = time.Second
		}
		return
	}

	if parser, ok := a.glossary[pattern.Origin()]; ok {
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
		reminder.AppendSubject(pattern)
	}

	a.buildReminder(speech, reminder)
}

func (a *Analyzer) getRandomStatement(seed int64) string {
	rand.Seed(seed)
	var statements = []string{
		`через 30 минут позвонить другу`,
		`завтра вынести мусор`,
		`вечером сделать домашнюю работу`,
		`в субботу купить корм коту`,
	}

	var rnd = rand.Intn(len(statements) - 1)

	return statements[rnd]
}
