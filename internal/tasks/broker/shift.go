package broker

import (
	"time"

	domain3 "github.com/sepuka/campaner/internal/command/domain"

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
		storedReminder *domain.Reminder
	)

	if storedReminder, err = b.reminderRepo.Get(reminder.ReminderId); err != nil {
		return err
	}

	if storedReminder.Whom != reminder.Whom {
		return errors2.NewWrongUserError(storedReminder.Whom, reminder.Whom)
	}

	if storedReminder.Status != domain.StatusNew {
		return errors2.NewWrongStatusError(storedReminder.Status, domain.StatusNew)
	}

	switch reminder.GetSubject() {
	case domain3.CommandOnTheEveId.String():
		reminder.When, err = b.onTheEve(storedReminder)
	case domain3.CommandBefore5Minutes.String():
		reminder.When, err = b.before5Minutes(storedReminder)
	default:
		return errors2.NewShiftError(0) //TODO tmp
	}

	if err != nil {
		return err
	}

	storedReminder.NotifyAt = time.Now().Add(reminder.When)

	return b.taskManager.Shift(storedReminder)
}

func (b *ShiftBroker) onTheEve(storedReminder *domain.Reminder) (time.Duration, error) {
	var (
		timeToEvent *calendar.Date
	)

	timeToEvent = calendar.NewDate(storedReminder.NotifyAt)
	if !calendar.IsNotSoon(timeToEvent.Until()) {
		return 0, errors2.NewShiftError(timeToEvent.Until())
	}

	timeToEvent = timeToEvent.Add(-calendar.Day).Evening()

	return timeToEvent.Until(), nil
}

func (b *ShiftBroker) before5Minutes(storedReminder *domain.Reminder) (time.Duration, error) {
	var (
		timeToEvent *calendar.Date
	)

	timeToEvent = calendar.NewDate(storedReminder.NotifyAt)
	if timeToEvent.Add(-5 * time.Minute).IsPast() {
		return 0, errors2.NewShiftError(timeToEvent.Until())
	}

	timeToEvent = timeToEvent.Add(-5 * time.Minute)

	return timeToEvent.Until(), nil
}
