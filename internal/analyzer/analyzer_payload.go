package analyzer

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/sepuka/campaner/internal/errors"

	domainApi "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	featureDomain "github.com/sepuka/campaner/internal/feature_toggling/domain"
	"go.uber.org/zap"
)

func (a *Analyzer) analyzePayload(msg context.Message, reminder *domain.Reminder) error {
	var (
		payload    domainApi.ButtonPayload
		err        error
		taskId     int64
		rawPayload = msg.Payload
		text       = domainApi.ButtonText(msg.Text)
	)

	if err = json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		a.logger.
			With(
				zap.String(`payload`, rawPayload),
				zap.Error(err),
			).
			Error(`analyze payload error`)
		return errors.NewInvalidSpeechPayloadFormatError(msg, err)
	}

	if taskId, err = strconv.ParseInt(payload.Button, 10, 64); err != nil {
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
	case domainApi.CancelButton:
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusCanceled
		// TODO снабдить reminder возможностью указывать кнопки, например  json keybord в новом поле бд
		// затем тут можно будет сделать так
		//reminder.Subject = []string{`напоминание отменено`}
		//reminder.When = time.Nanosecond
		// и юзер получит ответ без кнопок
	case domainApi.Later15MinButton:
		if !a.featureToggle.IsEnabled(reminder.Whom, featureDomain.Postpone) {
			return nil
		}

		reminder.Status = domain.StatusCopied
		reminder.ReminderId = int(taskId)
		reminder.When = time.Duration(15) * time.Minute

	}

	return nil
}
