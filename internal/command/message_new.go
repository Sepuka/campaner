package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sepuka/campaner/internal/analyzer"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/api"

	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/context"
)

type MessageNew struct {
	cfg      config.Server
	api      *api.SendMessage
	logger   *zap.SugaredLogger
	analyzer *analyzer.Analyzer
}

func NewMessageNew(
	cfg config.Server,
	api *api.SendMessage,
	logger *zap.SugaredLogger,
	analyzer *analyzer.Analyzer,
) *MessageNew {
	return &MessageNew{
		cfg:      cfg,
		api:      api,
		logger:   logger,
		analyzer: analyzer,
	}
}

func (obj *MessageNew) Exec(req *context.Request, resp http.ResponseWriter) error {
	var (
		output     = []byte(`ok`)
		text       = req.Object.Message.Text
		analyzeErr = fmt.Sprintf(`bad request: %s`, text)
		err        error
		reminder   = analyzer.NewReminder(
			int(req.Object.Message.PeerId),
		)
	)

	err = obj.analyzer.Analyze(text, reminder)
	if err != nil {
		reminder = analyzer.NewImmediateReminder(int(req.Object.Message.PeerId), analyzeErr)
		obj.
			logger.
			With(
				zap.String(`text`, text),
				zap.Error(err),
			).
			Error(`parse msg error`)
	}
	if reminder.When() == 0 {
		reminder = analyzer.NewImmediateReminder(int(req.Object.Message.PeerId), req.Object.Message.Text)
	}

	go obj.send(reminder)

	_, err = resp.Write(output)

	return err
}

func (obj *MessageNew) send(reminder *analyzer.Reminder) {
	var (
		err error
	)

	go obj.confirmMsg(reminder.When(), reminder.Whom())

	time.Sleep(reminder.When())

	if err = obj.api.Send(reminder.Whom(), reminder.What()); err != nil {
		obj.
			logger.
			With(
				zap.String(`text`, reminder.What()),
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
		now               = time.Now()
		todayMidnight     = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		yesterdayMidnight = todayMidnight.Add(time.Hour * 24)
	)

	if delay <= 3*time.Second {
		return
	}

	switch {
	case notificationTime.Before(todayMidnight):
		notifyTmpl := `напомню сегодня в %02d:%02d`
		text = fmt.Sprintf(notifyTmpl, notificationTime.Hour(), notificationTime.Minute())
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
