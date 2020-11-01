package broker

import "github.com/sepuka/campaner/internal/domain"

type NewBornBroker struct {
	reminderRepo domain.TaskManager
}

func NewNewBornBroker(repo domain.TaskManager) *NewBornBroker {
	return &NewBornBroker{reminderRepo: repo}
}

func (b *NewBornBroker) Service(reminder *domain.Reminder) error {
	return b.reminderRepo.Add(reminder)
}
