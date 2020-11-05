package broker

import (
	errors2 "github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/calendar"
	"github.com/sepuka/campaner/internal/domain"
)

type ShiftBroker struct {
	taskManager  domain.TaskManager
	reminderRepo domain.ReminderRepository
}

func NewShiftBroker(
	taskManager domain.TaskManager,
	repo domain.ReminderRepository,
) *ShiftBroker {
	return &ShiftBroker{
		taskManager:  taskManager,
		reminderRepo: repo,
	}
}

func (b *ShiftBroker) Service(reminder *domain.Reminder) error {
	var (
		err            error
		timeToEvent    *calendar.Date
		storedReminder *domain.Reminder
	)

	if storedReminder, err = b.reminderRepo.Get(reminder.ReminderId); err != nil {
		return err
	}

	if storedReminder.Whom != reminder.Whom {
		return errors2.NewWrongUserError(storedReminder.Whom, reminder.Whom)
	}

	timeToEvent = calendar.NewDate(storedReminder.NotifyAt)
	if !calendar.IsNotSoon(timeToEvent.Until()) {
		return errors2.NewShiftError(timeToEvent.Until())
	}

	if storedReminder.Status != domain.StatusNew {
		return errors2.NewWrongStatusError(storedReminder.Status, domain.StatusNew)
	}

	timeToEvent = timeToEvent.Add(-calendar.Day).Evening()
	reminder.When = timeToEvent.Until()

	return b.taskManager.Shift(storedReminder)
}
