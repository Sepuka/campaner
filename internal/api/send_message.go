package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strconv"
	"strings"

	"github.com/sepuka/campaner/internal/context"

	"github.com/sepuka/campaner/internal/config"
	"go.uber.org/zap"
)

const (
	apiVersion  = `5.103`
	apiEndpoint = `https://api.vk.com/method/`
	apiMethod   = `messages.send`
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

func (obj *SendMessage) Send(peerId int, text string) error {
	var (
		request      *http.Request
		response     *http.Response
		answer       = &context.Response{}
		dumpResponse []byte
		err          error
		params       = url2.Values{
			`v`:            []string{apiVersion},
			`access_token`: []string{obj.cfg.Api.Token},
			`message`:      []string{text},
			`peer_id`:      []string{strconv.Itoa(peerId)},
			`random_id`:    []string{strconv.FormatInt(rand.Int63(), 10)},
		}
		maskedAccessToken = fmt.Sprintf(`%s...`, obj.cfg.Api.Token[0:3])
		maskedParams      = strings.Replace(params.Encode(), obj.cfg.Api.Token, maskedAccessToken, 1)
		endpoint          = fmt.Sprintf(`%s/%s?%s`, apiEndpoint, apiMethod, params.Encode())
	)

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
