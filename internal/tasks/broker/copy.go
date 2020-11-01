package broker

import "github.com/sepuka/campaner/internal/domain"

type CopyBroker struct {
	reminderRepo domain.TaskManager
}

func NewCopyBroker(repo domain.TaskManager) *CopyBroker {
	return &CopyBroker{reminderRepo: repo}
}

func (b *CopyBroker) Service(reminder *domain.Reminder) error {
	return b.reminderRepo.Copy(reminder)
}
