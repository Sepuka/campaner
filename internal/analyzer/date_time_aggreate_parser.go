package analyzer

import (
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/speeches"
)

type DateTimeAggregateParser struct {
	parsers []Parser
}

func NewDateTimeAggregateParser(parsers []Parser) *DateTimeAggregateParser {
	return &DateTimeAggregateParser{
		parsers: parsers,
	}
}

func (p *DateTimeAggregateParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	const patternLength = 1
	var (
		err     error
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		return err
	}

	for _, parser := range p.parsers {
		if err = parser.Parse(speech, reminder); err == nil {
			return nil
		}
	}

	if err = speech.ApplyPattern(pattern); err != nil {
		return err
	}
	reminder.AppendSubject(pattern)

	return nil
}

func (p *DateTimeAggregateParser) Glossary() []string {
	return []string{
		`в`,
		`во`,
		`через`,
	}
}

func (p *DateTimeAggregateParser) PatternList() []string {
	return []string{
		`в среду`,
	}
}
