package broker

import "github.com/sepuka/campaner/internal/domain"

type BarrenBroker struct {
}

func NewBarrenBroker() *BarrenBroker {
	return &BarrenBroker{}
}

func (b *BarrenBroker) Service(reminder *domain.Reminder) error {
	return nil
}
