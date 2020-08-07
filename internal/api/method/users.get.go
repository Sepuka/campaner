package method

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strconv"
	"strings"

	"github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/config"
	"go.uber.org/zap"
)

type (
	UsersGet struct {
		cfg    *config.Config
		client *http.Client
		logger *zap.SugaredLogger
	}
)

func NewUsersGet(
	cfg *config.Config,
	client *http.Client,
	logger *zap.SugaredLogger,
) *UsersGet {
	return &UsersGet{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (obj *UsersGet) Send(peerId int) (*domain.User, error) {
	const (
		apiMethod = `users.get`
		fields    = `timezone`
	)
	var (
		request      *http.Request
		response     *http.Response
		userResponse = &domain.UserResponse{}
		dumpResponse []byte
		err          error
		params       = url2.Values{
			`v`:            []string{api.Version},
			`access_token`: []string{obj.cfg.Api.Token},
			`fields`:       []string{fields},
			`user_id`:      []string{strconv.Itoa(peerId)},
		}
		maskedAccessToken = fmt.Sprintf(`%s...`, obj.cfg.Api.Token[0:3])
		maskedParams      = strings.Replace(params.Encode(), obj.cfg.Api.Token, maskedAccessToken, 1)
		endpoint          = fmt.Sprintf(`%s/%s?%s`, api.Endpoint, apiMethod, params.Encode())
	)

	if request, err = http.NewRequest(`POST`, endpoint, nil); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`build api request error`)

		return nil, err
	}

	// TODO add an http short timeout
	if response, err = obj.client.Do(request); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`send api request error`)

		return nil, err
	}

	if dumpResponse, err = httputil.DumpResponse(response, true); err != nil {
		obj.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`dump api response error`)

		return nil, err
	}

	obj.
		logger.
		With(
			zap.String(`request`, maskedParams),
			zap.ByteString(`response`, dumpResponse),
		).
		Debug(`api message was sent`)

	if err = json.NewDecoder(response.Body).Decode(userResponse); err != nil {
		obj.
			logger.
			With(
				zap.Error(err),
				zap.ByteString(`response`, dumpResponse),
			).
			Error(`error while decoding api response`)

		return nil, err
	}

	return &userResponse.Response[0], nil
}
