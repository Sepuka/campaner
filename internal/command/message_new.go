package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/analyzer"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/api"

	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/context"
)

type MessageNew struct {
	cfg          config.Server
	api          *api.SendMessage
	logger       *zap.SugaredLogger
	analyzer     *analyzer.Analyzer
	reminderRepo domain.ReminderRepository
}

func NewMessageNew(
	cfg config.Server,
	api *api.SendMessage,
	logger *zap.SugaredLogger,
	analyzer *analyzer.Analyzer,
	repo domain.ReminderRepository,
) *MessageNew {
	return &MessageNew{
		cfg:          cfg,
		api:          api,
		logger:       logger,
		analyzer:     analyzer,
		reminderRepo: repo,
	}
}

func (obj *MessageNew) Exec(req *context.Request, resp http.ResponseWriter) error {
	var (
		err      error
		output   = []byte(`ok`)
		text     = req.Object.Message.Text
		reminder = domain.NewReminder(
			int(req.Object.Message.PeerId),
			text,
			time.Nanosecond,
		)
	)

	obj.analyzer.Analyze(text, reminder)
	if err = obj.reminderRepo.Persist(reminder); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
			).
			Error(`cannot save reminder`)
	}

	go obj.send(reminder)

	_, err = resp.Write(output)

	return err
}

func (obj *MessageNew) send(reminder *domain.Reminder) {
	var (
		err error
	)

	if !reminder.IsImmediate() {
		go obj.confirmMsg(reminder.When, reminder.Whom)
	}

	time.Sleep(reminder.When)

	if err = obj.api.Send(reminder.Whom, reminder.What); err != nil {
		obj.
			logger.
			With(
				zap.String(`text`, reminder.What),
				zap.Error(err),
			).
			Error(`send api message error`)
	}
}

func (obj *MessageNew) confirmMsg(delay time.Duration, whom int) {
	var (
		err  error
		text string

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

	if err = obj.api.Send(whom, text); err != nil {
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
