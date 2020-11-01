package tasks

import (
	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/domain"
)

type (
	Broker interface {
		Service(*domain.Reminder) error
	}
	BrokerMap map[int]Broker
)

type TaskBroker struct {
	brokerMap BrokerMap
}

func NewTaskBroker(brokerMap BrokerMap) *TaskBroker {
	return &TaskBroker{
		brokerMap: brokerMap,
	}
}

func (s *TaskBroker) Manage(reminder *domain.Reminder) error {
	var (
		broker Broker
		ok     bool
	)
	if broker, ok = s.brokerMap[reminder.Status]; ok {
		return broker.Service(reminder)
	}

	return errors.NewUnknownBrokerError(reminder.Status, `unknown status-key at broker map`)
}
