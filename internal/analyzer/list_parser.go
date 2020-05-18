package analyzer

import (
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"
)

const (
	limit                  = 5
	list           listCmd = `список`
	emptyTasksList         = `There aren't any tasks yet`
)

type (
	listCmd    string
	ListParser struct {
		reminderRepo domain.ReminderRepository
	}
)

func (lc listCmd) String() string {
	return string(lc)
}

func NewListParser(repo domain.ReminderRepository) *ListParser {
	return &ListParser{
		reminderRepo: repo,
	}
}

func (obj *ListParser) Parse(speech *speeches.Speech, reminder *domain.Reminder) error {
	var (
		err     error
		userId  = reminder.Whom
		models  []domain.Reminder
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(1); err != nil {
		return err
	}

	if models, err = obj.reminderRepo.Scheduled(userId, limit); err != nil {
		return err
	}

	if models == nil {
		reminder.What = emptyTasksList
	} else {
		var (
			schedule = make([]string, 0, len(models))
			list     string
		)
		for _, model := range models {
			schedule = append(schedule, model.String())
		}
		list = strings.Join(schedule, "\r\n")
		pattern = speeches.NewPattern([]string{list})
		reminder.AppendSubject(pattern)
	}

	reminder.When = time.Second

	return speech.ApplyPattern(pattern)
}

func (obj *ListParser) Glossary() []string {
	return []string{
		list.String(),
	}
}

func (obj *ListParser) PatternList() []string {
	return []string{
		`список`,
	}
}
