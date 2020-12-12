package broker

import (
	"strings"
	"testing"
	"time"

	domain2 "github.com/sepuka/campaner/internal/command/domain"

	"github.com/stretchr/testify/mock"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/errors"
	"github.com/stretchr/testify/assert"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/tasks"

	"github.com/go-pg/pg"

	"github.com/sepuka/campaner/internal/repository/mocks"
)

func TestShift(t *testing.T) {
	const (
		notExistentReminderId = 0
		ownerReminderId       = 1
		ownerWhom             = 1

		anotherReminderId = 2
		anotherWhom       = 2

		wrongStatusReminderId        = 3
		tooLateForEveReminderId      = 4
		tooLateFor5MinutesReminderId = 5
	)
	var (
		anotherReminder = &domain.Reminder{
			ReminderId: anotherReminderId,
			Whom:       anotherWhom,
			NotifyAt:   time.Time{},
			When:       0,
			Status:     0,
		}
		tooLateForEveReminder = &domain.Reminder{
			ReminderId: tooLateForEveReminderId,
			Whom:       ownerWhom,
			NotifyAt:   time.Now().Add(calendar.Day),
		}
		tooLateFor5MinReminder = &domain.Reminder{
			ReminderId: tooLateFor5MinutesReminderId,
			Whom:       ownerWhom,
			NotifyAt:   time.Now().Add(time.Minute),
		}
		wrongStatusReminder = &domain.Reminder{
			ReminderId: wrongStatusReminderId,
			Whom:       ownerWhom,
			NotifyAt:   time.Now().Add(calendar.Year),
			Status:     domain.StatusCanceled,
		}
		okReminder = &domain.Reminder{
			ReminderId: ownerReminderId,
			Whom:       ownerWhom,
			Status:     domain.StatusNew,
			NotifyAt:   time.Now().Add(10 * calendar.Day),
		}
		taskManager  = mocks.TaskManager{}
		reminderRepo = mocks.ReminderRepository{}
		broker       tasks.Broker
		actualError  error
		testCases    = map[string]struct {
			reminder *domain.Reminder
			wantErr  bool
			err      error
		}{
			`not existent reminder`: {
				reminder: &domain.Reminder{
					ReminderId: notExistentReminderId,
					Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
				},
				wantErr: true,
				err:     pg.ErrNoRows,
			},
			`stranger's reminder`: {
				reminder: &domain.Reminder{
					ReminderId: anotherReminderId,
					Whom:       ownerWhom,
					Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
				},
				wantErr: true,
				err:     errors.NewWrongUserError(ownerWhom, anotherWhom),
			},
			`it is too late (on the eve)`: {
				reminder: &domain.Reminder{
					ReminderId: tooLateForEveReminderId,
					Whom:       ownerWhom,
					Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
				},
				wantErr: true,
				err:     errors.NewShiftError(10 * time.Second),
			},
			`it is too late (before 5 minutes)`: {
				reminder: &domain.Reminder{
					ReminderId: tooLateFor5MinutesReminderId,
					Whom:       ownerWhom,
					Subject:    strings.Split(domain2.CommandBefore5Minutes.String(), ` `),
				},
				wantErr: true,
				err:     errors.NewShiftError(10 * time.Second),
			},
			`wrong status`: {
				reminder: &domain.Reminder{
					ReminderId: wrongStatusReminderId,
					Whom:       ownerWhom,
					Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
				},
				wantErr: true,
				err:     errors.NewWrongStatusError(wrongStatusReminder.Status, domain.StatusNew),
			},
			`ok`: {
				reminder: &domain.Reminder{
					ReminderId: ownerReminderId,
					Whom:       ownerWhom,
					Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
				},
				wantErr: false,
			},
		}
	)

	reminderRepo.
		On(`Get`, notExistentReminderId).Return(nil, pg.ErrNoRows).
		On(`Get`, tooLateForEveReminderId).Return(tooLateForEveReminder, nil).
		On(`Get`, tooLateFor5MinutesReminderId).Return(tooLateFor5MinReminder, nil).
		On(`Get`, anotherReminderId).Return(anotherReminder, nil).
		On(`Get`, wrongStatusReminderId).Return(wrongStatusReminder, nil).
		On(`Get`, ownerReminderId).Return(okReminder, nil)
	taskManager.
		On(`Shift`, tooLateForEveReminder).Times(0).
		On(`Shift`, tooLateFor5MinReminder).Times(0).
		On(`Shift`, anotherReminder).Times(0).
		On(`Shift`, wrongStatusReminder).Times(0).
		On(`Shift`, okReminder, mock.Anything).Once().Return(nil)

	for testName, testCase := range testCases {
		broker = NewShiftBroker(taskManager, reminderRepo)
		actualError = broker.Service(testCase.reminder)
		if testCase.wantErr != (actualError != nil) {
			t.Errorf(`unexpected error %v in %s`, actualError, testName)
			return
		}
		if actualError != nil {
			assert.EqualError(t, actualError, testCase.err.Error())
		}
	}
}

func TestShift_TimeCheck_OnTheEve(t *testing.T) {
	var (
		notifyAt         = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.Local).Add(10 * calendar.Day)
		expectedNotifyAt = notifyAt.Add(-15 * time.Hour)
		taskManager      = mocks.TaskManager{}
		reminderRepo     = mocks.ReminderRepository{}
		storedReminder   = &domain.Reminder{
			ReminderId: 1,
			Whom:       2,
			NotifyAt:   notifyAt,
			Status:     domain.StatusNew,
		}
		expectedReminder = storedReminder
		flowReminder     = &domain.Reminder{
			ReminderId: 1,
			Whom:       2,
			Subject:    strings.Split(domain2.CommandOnTheEveId.String(), ` `),
		}
		broker tasks.Broker
		err    error
	)

	expectedReminder.NotifyAt = expectedNotifyAt
	reminderRepo.On(`Get`, mock.Anything).Return(storedReminder, nil)
	taskManager.On(`Shift`, expectedReminder).Once().Return(nil)
	broker = NewShiftBroker(taskManager, reminderRepo)
	err = broker.Service(flowReminder)

	assert.NoError(t, err)
}

func TestShift_TimeCheck_Before5Min(t *testing.T) {
	var (
		notifyAtTomorrowMorning = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.Local).Add(calendar.Day)
		expectedNotifyAt        = notifyAtTomorrowMorning.Add(-5 * time.Minute)
		taskManager             = mocks.TaskManager{}
		reminderRepo            = mocks.ReminderRepository{}
		storedReminder          = &domain.Reminder{
			ReminderId: 1,
			Whom:       2,
			NotifyAt:   notifyAtTomorrowMorning,
			Status:     domain.StatusNew,
		}
		expectedReminder = storedReminder
		flowReminder     = &domain.Reminder{
			ReminderId: 1,
			Whom:       2,
			Subject:    strings.Split(domain2.CommandBefore5Minutes.String(), ` `),
		}
		broker tasks.Broker
		err    error
	)

	expectedReminder.NotifyAt = expectedNotifyAt
	reminderRepo.On(`Get`, mock.Anything).Return(storedReminder, nil)
	taskManager.On(`Shift`, expectedReminder).Once().Return(nil)
	broker = NewShiftBroker(taskManager, reminderRepo)
	err = broker.Service(flowReminder)

	assert.NoError(t, err)
}
