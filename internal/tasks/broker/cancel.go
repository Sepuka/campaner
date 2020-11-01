package broker

import "github.com/sepuka/campaner/internal/domain"

type CancelBroker struct {
	reminderRepo domain.TaskManager
}

func NewCancelBroker(repo domain.TaskManager) *CancelBroker {
	return &CancelBroker{reminderRepo: repo}
}

func (b *CancelBroker) Service(reminder *domain.Reminder) error {
	return b.reminderRepo.Cancel(reminder)
}
