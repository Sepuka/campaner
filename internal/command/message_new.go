package command

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/sepuka/campaner/internal/tasks"

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
	api      *method.SendMessage
	logger   *zap.SugaredLogger
	analyzer *analyzer.Analyzer
	store    *tasks.TaskBroker
}

func NewMessageNew(
	api *method.SendMessage,
	logger *zap.SugaredLogger,
	analyzer *analyzer.Analyzer,
	store *tasks.TaskBroker,
) *MessageNew {
	return &MessageNew{
		api:      api,
		logger:   logger,
		analyzer: analyzer,
		store:    store,
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
		if errors.IsNotATimeError(err) {
			var randomSubject = fmt.Sprintf(`Попробуйте фразу: "%s"`, obj.getRandomStatement(time.Now().Unix()))
			reminder.RewriteSubject(randomSubject)

			return obj.api.SendFlat(reminder.Whom, reminder.GetSubject())
		}
		return err
	}

	if err = obj.store.Manage(reminder); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
				zap.Int(`reminder_id`, reminder.ReminderId),
				zap.Int(`user_id`, reminder.Whom),
				zap.Int(`status`, reminder.Status),
			).
			Error(`cannot manage reminder`)
		return err
	}

	go obj.confirmMsg(reminder)

	return err
}

func (obj *MessageNew) confirmMsg(reminder *domain.Reminder) {
	var (
		err              error
		text, notifyTmpl string
		whom             = reminder.Whom

		notificationTime  = time.Now().Add(reminder.When)
		todayMidnight     = calendar.NextMidnight()
		yesterdayMidnight = calendar.LastMidnight()
	)

	if reminder.IsCancelled() {
		text = `напоминание отменено`
	} else if reminder.IsBarren() {
		text = reminder.GetSubject()
	} else if reminder.IsImmediate() == false {
		switch {
		case notificationTime.Before(todayMidnight):
			notifyTmpl = `напомню сегодня в %02d:%02d:%02d`
			text = fmt.Sprintf(notifyTmpl, notificationTime.Hour(), notificationTime.Minute(), notificationTime.Second())
		case notificationTime.Before(yesterdayMidnight):
			notifyTmpl = `напомню завтра в %02d:%02d`
			text = fmt.Sprintf(notifyTmpl, notificationTime.Hour(), notificationTime.Minute())
		default:
			notifyTmpl = `напомню об этом %d.%02d в %02d:%02d`
			text = fmt.Sprintf(notifyTmpl, notificationTime.Day(), notificationTime.Month(), notificationTime.Hour(), notificationTime.Minute())
		}
	} else {
		return
	}

	if err = obj.api.SendIntention(whom, text, reminder); err != nil {
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

func (obj *MessageNew) getRandomStatement(seed int64) string {
	rand.Seed(seed)
	var statements = []string{
		`через 30 минут позвонить другу`,
		`завтра вынести мусор`,
		`вечером сделать домашнюю работу`,
		`в субботу купить корм коту`,
	}

	var rnd = rand.Intn(len(statements) - 1)

	return statements[rnd]
}
