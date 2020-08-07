package command

import (
	"fmt"
	"net/http"
	"time"

	apiDomain "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/api/method"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/analyzer"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/context"
)

type MessageNew struct {
	api          *method.MessagesSend
	apiUsersGet  *method.UsersGet
	logger       *zap.SugaredLogger
	analyzer     *analyzer.Analyzer
	reminderRepo domain.ReminderRepository
}

// TODO: encapsulate API to one arg
func NewMessageNew(
	api *method.MessagesSend,
	apiUsersGet *method.UsersGet,
	logger *zap.SugaredLogger,
	analyzer *analyzer.Analyzer,
	repo domain.ReminderRepository,
) *MessageNew {
	return &MessageNew{
		api:          api,
		apiUsersGet:  apiUsersGet,
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
		reminder = domain.NewReminder(int(req.Object.Message.PeerId))
		reqCtx   *context.Context
	)

	if reqCtx, err = obj.buildContext(req); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
			).
			Error(`build context error`)
	}

	obj.analyzer.Analyze(text, reminder, reqCtx)
	if err = obj.reminderRepo.Add(reminder); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
			).
			Error(`cannot save reminder`)
	}

	if !reminder.IsImmediate() {
		go obj.confirmMsg(reminder.When, reminder.Whom)
	}

	_, err = resp.Write(output)

	return err
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

func (obj *MessageNew) buildContext(req *context.Request) (*context.Context, error) {
	var (
		err    error
		userId = int(req.Object.Message.PeerId)
		user   *apiDomain.User
	)

	if user, err = obj.apiUsersGet.Send(userId); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
			).
			Error(`send api message error (users.get)`)
	}

	return &context.Context{
		User: &context.User{Timezone: user.TimeZone},
	}, nil
}

func (obj *MessageNew) Precept() []string {
	return []string{
		`message_new`,
	}
}
