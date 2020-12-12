package analyzer

import (
	"encoding/json"

	"github.com/sepuka/campaner/internal/analyzer/payload"

	"github.com/sepuka/campaner/internal/errors"

	domainApi "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	"go.uber.org/zap"
)

func (a *Analyzer) analyzePayload(msg context.Message, reminder *domain.Reminder) error {
	var (
		err         error
		rawPayload  = msg.Payload
		payloadData domainApi.ButtonPayload
	)

	if err = json.Unmarshal([]byte(rawPayload), &payloadData); err != nil {
		return a.buildPayloadError(msg, reminder.Whom, err, `invalid JSON`)
	}

	switch {
	case payloadData.IsStartButton():
		if err = payload.HandleStartButtonPayload(payloadData, reminder); err != nil {
			return a.buildPayloadError(msg, reminder.Whom, err, `something wrong with start button`)
		}
	case payloadData.IsChangeTaskButton():
		if err = payload.HandleChangeTaskPayload(msg.Text, payloadData, reminder); err != nil {
			return a.buildPayloadError(msg, reminder.Whom, err, `cannot parse task_id`)
		}
	case payloadData.IsOKButton():
		if err = payload.HandleOKButtonPayload(msg.Text, payloadData, reminder); err != nil {
			return a.buildPayloadError(msg, reminder.Whom, err, `cannot handle OK button`)
		}
	default:
		return a.buildPayloadError(msg, reminder.Whom, err, `got unknown button payload`)
	}

	return nil
}

func (a *Analyzer) buildPayloadError(msg context.Message, userId int, err error, text string) error {
	var (
		rawPayload = msg.Payload
	)

	a.
		logger.
		With(
			zap.String(`json`, rawPayload),
			zap.Int(`user_id`, userId),
			zap.Error(err),
		).
		Error(text)

	return errors.NewInvalidSpeechPayloadButtonError(msg, err)
}
