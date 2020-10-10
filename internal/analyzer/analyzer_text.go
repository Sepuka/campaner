package analyzer

import (
	"fmt"
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/errors"
	"github.com/sepuka/campaner/internal/speeches"
)

func (a *Analyzer) analyzeText(text string, reminder *domain.Reminder) error {
	return a.buildReminder(speeches.NewSpeech(text), reminder)
}

func (a *Analyzer) buildReminder(speech *speeches.Speech, reminder *domain.Reminder) error {
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
			var randomSubject = fmt.Sprintf(`Попробуйте фразу: "%s"`, a.getRandomStatement(time.Now().Unix()))
			reminder.RewriteSubject(randomSubject)
			reminder.When = time.Second
		}
		return nil
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
			return nil
		}
	} else {
		if err = speech.ApplyPattern(pattern); err != nil {
			return nil
		}
		reminder.AppendSubject(pattern)
	}

	return a.buildReminder(speech, reminder)
}
