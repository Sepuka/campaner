package broker

import (
	"testing"
	"time"

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

		wrongStatusReminderId = 3
		tooLateReminderId     = 4
	)
	var (
		anotherReminder = &domain.Reminder{
			ReminderId: anotherReminderId,
			Whom:       anotherWhom,
			NotifyAt:   time.Time{},
			When:       0,
			Status:     0,
		}
		tooLateReminder = &domain.Reminder{
			ReminderId: tooLateReminderId,
			Whom:       ownerWhom,
			NotifyAt:   time.Now().Add(calendar.Day),
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
				},
				wantErr: true,
				err:     pg.ErrNoRows,
			},
			`stranger's reminder`: {
				reminder: &domain.Reminder{
					ReminderId: anotherReminderId,
					Whom:       ownerWhom,
				},
				wantErr: true,
				err:     errors.NewWrongUserError(ownerWhom, anotherWhom),
			},
			`it is too late`: {
				reminder: &domain.Reminder{
					ReminderId: tooLateReminderId,
					Whom:       ownerWhom,
				},
				wantErr: true,
				err:     errors.NewShiftError(10 * time.Second),
			},
			`wrong status`: {
				reminder: &domain.Reminder{
					ReminderId: wrongStatusReminderId,
					Whom:       ownerWhom,
				},
				wantErr: true,
				err:     errors.NewWrongStatusError(wrongStatusReminder.Status, domain.StatusNew),
			},
			`ok`: {
				reminder: &domain.Reminder{
					ReminderId: ownerReminderId,
					Whom:       ownerWhom,
				},
				wantErr: false,
			},
		}
	)

	reminderRepo.
		On(`Get`, notExistentReminderId).Return(nil, pg.ErrNoRows).
		On(`Get`, tooLateReminderId).Return(tooLateReminder, nil).
		On(`Get`, anotherReminderId).Return(anotherReminder, nil).
		On(`Get`, wrongStatusReminderId).Return(wrongStatusReminder, nil).
		On(`Get`, ownerReminderId).Return(okReminder, nil)
	taskManager.
		On(`Shift`, tooLateReminder).Times(0).
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

func TestShift_TimeChecks(t *testing.T) {
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
