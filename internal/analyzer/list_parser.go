package analyzer

import (
	"strings"

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
		err    error
		userId = reminder.Whom
		models []domain.Reminder
	)

	if models, err = obj.reminderRepo.Scheduled(userId, limit); err != nil {
		return err
	}

	if models == nil {
		reminder.What = emptyTasksList
	} else {
		var schedule = make([]string, 0, len(models))
		for _, m := range models {
			schedule = append(schedule, m.String())
		}
		reminder.What = strings.Join(schedule, "\r\n")
	}

	return err
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
