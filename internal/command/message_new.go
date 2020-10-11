package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sepuka/campaner/internal/errors"

	"github.com/sepuka/campaner/internal/api"

	"github.com/sepuka/campaner/internal/api/method"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/analyzer"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/context"
)

type MessageNew struct {
	api          *method.SendMessage
	logger       *zap.SugaredLogger
	analyzer     *analyzer.Analyzer
	reminderRepo domain.TaskManager
}

func NewMessageNew(
	api *method.SendMessage,
	logger *zap.SugaredLogger,
	analyzer *analyzer.Analyzer,
	repo domain.TaskManager,
) *MessageNew {
	return &MessageNew{
		api:          api,
		logger:       logger,
		analyzer:     analyzer,
		reminderRepo: repo,
	}
}

func (obj *MessageNew) Exec(req *context.Request, resp http.ResponseWriter) error {
	var (
		err      error
		reminder = domain.NewReminder(int(req.Object.Message.PeerId))
		msg      = req.Object.Message
	)

	if _, err = resp.Write(api.Response()); err != nil {
		return err
	}

	if err = obj.analyzer.Analyze(msg, reminder); err != nil {
		return err
	}

	switch reminder.Status {
	case domain.StatusNew:
		if err = obj.reminderRepo.Add(reminder); err != nil {
			obj.
				logger.
				With(
					zap.Error(err),
				).
				Error(`cannot save reminder`)
		}
	case domain.StatusCanceled:
		if err = obj.reminderRepo.Cancel(reminder); err != nil {
			obj.
				logger.
				With(
					zap.Error(err),
					zap.Int(`task_id`, reminder.ReminderId),
					zap.Int(`user_id`, reminder.Whom),
				).
				Error(`cannot save reminder`)
		}
	case domain.StatusCopied:
		if err = obj.reminderRepo.Copy(reminder); err != nil {
			obj.
				logger.
				With(
					zap.Int(`task_id`, reminder.ReminderId),
					zap.Int(`user_id`, reminder.Whom),
					zap.Error(err),
				).
				Error(`cannot prolong task`)
			return errors.NewStorageError(`taskManager`, err)
		}
	}

	if !reminder.IsImmediate() {
		go obj.confirmMsg(reminder)
	}

	return err
}

func (obj *MessageNew) confirmMsg(reminder *domain.Reminder) {
	var (
		err   error
		text  string
		delay = reminder.When
		whom  = reminder.Whom

		notificationTime  = time.Now().Add(delay)
		todayMidnight     = calendar.NextMidnight()
		yesterdayMidnight = calendar.LastMidnight()
	)

	switch {
	case notificationTime.Before(todayMidnight):
		notifyTmpl := `напомню сегодня в %02d:%02d:%02d`
		text = fmt.Sprintf(notifyTmpl, notificationTime.Hour(), notificationTime.Minute(), notificationTime.Second())
	case notificationTime.Before(yesterdayMidnight):
		notifyTmpl := `напомню завтра в %02d:%02d`
		text = fmt.Sprintf(notifyTmpl, notificationTime.Hour(), notificationTime.Minute())
	default:
		notifyTmpl := `напомню об этом %d.%02d в %02d:%02d`
		text = fmt.Sprintf(notifyTmpl, notificationTime.Day(), notificationTime.Month(), notificationTime.Hour(), notificationTime.Minute())
	}

	if err = obj.api.SendIntention(whom, text, reminder.ReminderId); err != nil {
		obj.
			logger.
			With(
				zap.String(`text`, text),
				zap.Error(err),
			).
			Error(`send api message error (confirmation)`)
	}
}

func (obj *MessageNew) Precept() []string {
	return []string{
		`message_new`,
	}
}
