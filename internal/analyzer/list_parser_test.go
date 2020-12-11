package analyzer

import (
	"testing"
	"time"

	"github.com/sepuka/campaner/internal/errors"

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
			{
				What:     `The first scheduled notification`,
				When:     time.Second,
				NotifyAt: time.Date(1984, 8, 31, 0, 0, 0, 0, time.Local),
			},
			{
				What:     `The second scheduled notification`,
				When:     time.Second,
				NotifyAt: time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local),
			},
		}
		actualReminder = domain.NewImmediateReminder(0, ``)
		expectedText   = "\"The first scheduled notification\" at 1984-08-31 00:00:00\r\n" +
			"\"The second scheduled notification\" at 2000-12-31 23:59:59"
		err    error
		speech = speeches.NewSpeech(list.String())
	)

	repo.On(`Scheduled`, mock.Anything, mock.Anything).Return(reminders, nil)

	var parser = NewListParser(repo)
	err = parser.Parse(speech, actualReminder)
	assert.Equal(t, expectedText, actualReminder.GetSubject())
	assert.NoError(t, err)

	patternWasApplied(t, speech)
}

func TestListNotificationsAnalyzer_noTasks(t *testing.T) {
	var (
		repo           = mocks.ReminderRepository{}
		reminders      []domain.Reminder
		actualReminder = domain.NewImmediateReminder(0, ``)
		expectedText   = `There aren't any tasks yet`
		err            error
		speech         = speeches.NewSpeech(list.String())
	)

	repo.On(`Scheduled`, mock.Anything, mock.Anything).Return(reminders, nil)

	var parser = NewListParser(repo)
	err = parser.Parse(speech, actualReminder)
	assert.Equal(t, expectedText, actualReminder.GetSubject())
	assert.NoError(t, err)

	patternWasApplied(t, speech)
}

func patternWasApplied(t *testing.T, speech *speeches.Speech) {
	var (
		err     error
		pattern *speeches.Pattern
	)
	pattern, err = speech.TryPattern(1)
	assert.Nil(t, pattern, `pattern must be empty after applying`)
	assert.EqualError(t, err, errors.NewSpeechIsOverError(1, list.String()).Error())
}
