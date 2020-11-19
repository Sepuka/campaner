package analyzer

import (
	"time"

	"github.com/sepuka/campaner/internal/api/method"

	payload2 "github.com/sepuka/campaner/internal/analyzer/payload"

	"github.com/sepuka/campaner/internal/errors"

	domainApi "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	"go.uber.org/zap"
)

func (a *Analyzer) analyzePayload(msg context.Message, reminder *domain.Reminder) error {
	var (
		err        error
		taskId     int64
		rawPayload = msg.Payload
		text       = domainApi.ButtonText(msg.Text)
	)

	if taskId, err = payload2.GetTaskId(rawPayload); err != nil {
		a.
			logger.
			With(
				zap.String(`json`, rawPayload),
				zap.Int(`user_id`, reminder.Whom),
				zap.Error(err),
			).
			Error(`cannot parse task_id`)
		return errors.NewInvalidSpeechPayloadButtonError(msg, err)
	}

	switch text {
	case method.CancelButton:
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusCanceled
		// TODO снабдить reminder возможностью указывать кнопки, например  json keybord в новом поле бд
		// затем тут можно будет сделать так
		//reminder.Subject = []string{`напоминание отменено`}
		//reminder.When = time.Nanosecond
		// и юзер получит ответ без кнопок
	case method.Later15MinButton:
		reminder.Status = domain.StatusCopied
		reminder.ReminderId = int(taskId)
		reminder.When = time.Duration(15) * time.Minute
	case method.OKButton:
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusBarren
	case method.OnTheEve, method.Before5Minutes:
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusShifted
		reminder.RewriteSubject(msg.Text)
	}

	return nil
}
