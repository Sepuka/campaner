package analyzer

import (
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/repository/mocks"

	"github.com/sepuka/campaner/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListNotificationsAnalyzer(t *testing.T) {
	var (
		repo      = mocks.ReminderRepository{}
		reminders = []domain.Reminder{
			*domain.NewReminder(0, `The first scheduled notification`, time.Second),
			*domain.NewReminder(0, `The second scheduled notification`, time.Second),
		}
		actualReminder = domain.NewImmediateReminder(0, ``)
		expectedText   = "\"The first scheduled notification\" at 1984-08-31 00:00:00\r\n" +
			"\"The second scheduled notification\" at 2000-12-31 23:59:59"
		err error
	)

	reminders[0].NotifyAt = time.Date(1984, 8, 31, 0, 0, 0, 0, time.Local)
	reminders[1].NotifyAt = time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
	repo.On(`Scheduled`, mock.Anything, mock.Anything).Return(reminders, nil)

	var parser = NewListParser(repo)
	err = parser.Parse(speeches.NewSpeech(`список`), actualReminder)
	assert.Equal(t, expectedText, actualReminder.What)
	assert.NoError(t, err)
}

func TestListNotificationsAnalyzer_noTasks(t *testing.T) {
	var (
		repo           = mocks.ReminderRepository{}
		reminders      []domain.Reminder
		actualReminder = domain.NewImmediateReminder(0, ``)
		expectedText   = `There aren't any tasks yet`
		err            error
	)

	repo.On(`Scheduled`, mock.Anything, mock.Anything).Return(reminders, nil)

	var parser = NewListParser(repo)
	err = parser.Parse(speeches.NewSpeech(`список`), actualReminder)
	assert.Equal(t, expectedText, actualReminder.What)
	assert.NoError(t, err)
}
