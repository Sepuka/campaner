package method

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strings"

	"github.com/sepuka/campaner/internal/calendar"
	domain2 "github.com/sepuka/campaner/internal/domain"

	"github.com/google/go-querystring/query"
	"github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/api/domain"
	featureDomain "github.com/sepuka/campaner/internal/feature_toggling/domain"

	"github.com/sepuka/campaner/internal/context"

	"github.com/sepuka/campaner/internal/config"
	"go.uber.org/zap"
)

const (
	apiMethod = `messages.send`
)

type (
	SendMessage struct {
		cfg           *config.Config
		client        api.HTTPClient
		logger        *zap.SugaredLogger
		featureToggle featureDomain.FeatureToggle
	}
)

func NewSendMessage(
	cfg *config.Config,
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	feature featureDomain.FeatureToggle,
) *SendMessage {
	return &SendMessage{
		cfg:           cfg,
		client:        client,
		logger:        logger,
		featureToggle: feature,
	}
}

func (obj *SendMessage) SendIntention(peerId int, text string, reminder *domain2.Reminder) error {
	var (
		err      error
		params   url2.Values
		keyboard = domain.Keyboard{
			OneTime: true,
			Buttons: cancel(reminder.ReminderId),
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

	if reminder.IsShifted() {
		keyboard.Buttons = cancel(reminder.ReminderId)
	} else if calendar.IsNotSoon(reminder.When) {
		keyboard.Buttons = cancelWithEve(reminder.ReminderId)
	} else if obj.hasNotAnyButtons(reminder) {
		keyboard.Buttons = [][]domain.Button{}
	}

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

func (obj *SendMessage) SendNotification(peerId int, text string, remindId int) error {
	var (
		payload = domain.MessagesSend{
			Message:     text,
			AccessToken: obj.cfg.Api.Token,
			ApiVersion:  api.Version,
			PeerId:      peerId,
			RandomId:    api.Rnd(),
		}
		err      error
		js       []byte
		keyboard = domain.Keyboard{
			OneTime: true,
			Buttons: delayAndOk(remindId),
		}
	)

	if js, err = json.Marshal(keyboard); err != nil {
		obj.
			logger.
			With(
				zap.Any(`request`, keyboard),
				zap.Error(err),
			).
			Errorf(`build keyboard query string error`)

		return err
	}

	payload.Keyboard = string(js)

	return obj.send(payload)
}

func (obj *SendMessage) SendFlat(peerId int, text string) error {
	var (
		payload = domain.MessagesSend{
			Message:     text,
			AccessToken: obj.cfg.Api.Token,
			ApiVersion:  api.Version,
			PeerId:      peerId,
			RandomId:    api.Rnd(),
		}
	)

	return obj.send(payload)
}

func (obj *SendMessage) send(msgStruct domain.MessagesSend) error {
	var (
		request      *http.Request
		response     *http.Response
		answer       = &context.Response{}
		dumpResponse []byte
		err          error
		params       url2.Values
		maskedParams string
		endpoint     string
	)

	if params, err = query.Values(msgStruct); err != nil {
		obj.
			logger.
			With(zap.Error(err)).
			Errorf(`build request query string error`)

		return err
	}

	endpoint = fmt.Sprintf(`%s/%s?%s`, api.Endpoint, apiMethod, params.Encode())
	maskedParams = obj.cfg.Api.MaskedToken(endpoint)

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
		Info(`api message sent`)

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

func (obj *SendMessage) hasNotAnyButtons(reminder *domain2.Reminder) bool {
	// TODO add Status button
	return reminder.IsCancelled() || reminder.IsBarren()
}
