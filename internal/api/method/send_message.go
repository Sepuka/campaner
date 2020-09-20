package method

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"

	"github.com/sepuka/campaner/internal/config"
	"go.uber.org/zap"
)

const (
	apiMethod = `messages.send`
)

type (
	SendMessage struct {
		cfg    *config.Config
		client *http.Client
		logger *zap.SugaredLogger
	}
)

func NewSendMessage(
	cfg *config.Config,
	client *http.Client,
	logger *zap.SugaredLogger,
) *SendMessage {
	return &SendMessage{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (obj *SendMessage) SendIntention(peerId int, text string, cancelId int) error {
	var (
		err      error
		params   url2.Values
		keyboard = domain.Keyboard{
			OneTime: true,
			Buttons: [][]domain.Button{
				{
					{
						Color: `negative`,
						Action: domain.Action{
							Type:    domain.TextButtonType,
							Label:   domain.CancelButton,
							Payload: domain.ButtonPayload{Button: fmt.Sprintf(`%d`, cancelId)}.String(),
						},
					},
				},
			},
		}

		payload = domain.MessagesSend{
			Message:     text,
			AccessToken: obj.cfg.Api.Token,
			ApiVersion:  api.Version,
			PeerId:      peerId,
			RandomId:    rand.Int63(),
		}

		maskedAccessToken = fmt.Sprintf(`%s...`, obj.cfg.Api.Token[0:3])
		maskedParams      = strings.Replace(params.Encode(), obj.cfg.Api.Token, maskedAccessToken, 1)
		js                []byte
	)

	if js, err = json.Marshal(keyboard); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`build keyboard query string error`)

		return err
	}

	payload.Keyboard = string(js)

	return obj.send(payload)
}

func (obj *SendMessage) SendNotification(peerId int, text string) error {
	var (
		payload = domain.MessagesSend{
			Message:     text,
			AccessToken: obj.cfg.Api.Token,
			ApiVersion:  api.Version,
			PeerId:      peerId,
			RandomId:    rand.Int63(),
		}
	)

	return obj.send(payload)
}

func (obj *SendMessage) send(queryArgs domain.MessagesSend) error {
	var (
		request      *http.Request
		response     *http.Response
		answer       = &context.Response{}
		dumpResponse []byte
		err          error
		params       url2.Values

		maskedAccessToken = fmt.Sprintf(`%s...`, obj.cfg.Api.Token[0:3])
		maskedParams      = strings.Replace(params.Encode(), obj.cfg.Api.Token, maskedAccessToken, 1)
		endpoint          string
	)

	if params, err = query.Values(queryArgs); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`build request query string error`)

		return err
	}

	endpoint = fmt.Sprintf(`%s/%s?%s`, api.Endpoint, apiMethod, params.Encode())

	if request, err = http.NewRequest(`POST`, endpoint, nil); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`build api request error`)

		return err
	}

	if response, err = obj.client.Do(request); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`send api request error`)

		return err
	}

	if dumpResponse, err = httputil.DumpResponse(response, true); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`dump api response error`)

		return err
	}

	obj.
		logger.
		With(
			zap.String(`request`, maskedParams),
			zap.ByteString(`response`, dumpResponse),
		).
		Debug(`api message was sent`)

	if err = json.NewDecoder(response.Body).Decode(answer); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
				zap.ByteString(`response`, dumpResponse),
			).
			Error(`error while decoding api response`)

		return err
	}

	if len(answer.Error.Message) > 0 {
		obj.
			logger.
			With(
				zap.Int32(`code`, answer.Error.Code),
				zap.String(`message`, answer.Error.Message),
			).
			Error(`failed api answer`)
	}

	return nil
}
