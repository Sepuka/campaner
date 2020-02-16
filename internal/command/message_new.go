package command

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/api"

	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/context"
)

type MessageNew struct {
	cfg    config.Server
	api    *api.SendMessage
	logger *zap.SugaredLogger
}

func NewMessageNew(
	cfg config.Server,
	api *api.SendMessage,
	logger *zap.SugaredLogger,
) *MessageNew {
	return &MessageNew{
		cfg:    cfg,
		api:    api,
		logger: logger,
	}
}

func (obj *MessageNew) Exec(req *context.Request, resp http.ResponseWriter) error {
	var (
		output = []byte(`ok`)
		text   = req.Object.Message.Text
		err    error
	)

	if err = obj.api.Send(int(req.Object.Message.PeerId), text); err != nil {
		obj.
			logger.
			With(
				zap.String(`text`, text),
				zap.Error(err),
			).
			Error(`send api message error`)
	}

	_, err = resp.Write(output)

	return err
}

func (obj *MessageNew) Precept() []string {
	return []string{
		`message_new`,
	}
}
