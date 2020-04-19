package analyzer

import (
	"strings"

	"github.com/sepuka/campaner/internal/domain"
)

const (
	limit         = 5
	list  listCmd = `список`
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

func (obj *ListParser) Parse(words []string, reminder *domain.Reminder) ([]string, error) {
	var (
		err    error
		userId = reminder.Whom
		models []domain.Reminder
	)

	if models, err = obj.reminderRepo.Scheduled(userId, limit); err != nil {
		return words, err
	}

	var schedule = make([]string, 0, len(models))
	for _, m := range models {
		schedule = append(schedule, m.String())
	}

	reminder.What = strings.Join(schedule, "\r\n")

	return words[1:], err
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
