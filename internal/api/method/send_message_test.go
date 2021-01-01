package method

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/sepuka/campaner/internal/api"

	domain3 "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"
	domain2 "github.com/sepuka/campaner/internal/domain"

	mocks2 "github.com/sepuka/campaner/internal/api/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/sepuka/campaner/internal/feature_toggling/domain"

	"github.com/sepuka/campaner/internal/feature_toggling/toggle/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/config"
)

func TestSendNotification_withPostponeButtons(t *testing.T) {
	var (
		cfg = &config.Config{
			Api: config.Api{
				Token: `some_token`,
			},
		}
		client    = &mocks2.HTTPClient{}
		rnd       = domain3.NewRnder()
		answer, _ = json.Marshal(&context.Response{})
		response  = &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(answer)),
		}
		logger = zap.NewNop().Sugar()
		ft     = mocks.FeatureToggle{}
		err    error
		peerId = 1
	)
	ft.On(`IsEnabled`, peerId, domain.Postpone).Return(true)
	client.On(`Do`, mock.Anything).Return(response, nil)

	var sm = NewSendMessage(cfg, client, logger, ft, rnd)
	err = sm.SendNotification(peerId, ``, 1)
	assert.NoError(t, err)
}

func TestSendIntention_Buttons(t *testing.T) {
	const (
		rndId = int64(1)
	)
	var (
		peerId = 1
		cfg    = &config.Config{
			Api: config.Api{
				Token: `some_token`,
			},
		}
		client    = &mocks2.HTTPClient{}
		logger    = zap.NewNop().Sugar()
		ft        = mocks.FeatureToggle{}
		rnd       = &mocks2.Rnder{}
		sm        = NewSendMessage(cfg, client, logger, ft, rnd)
		answer, _ = json.Marshal(&context.Response{})
		response  *http.Response
		request   *http.Request
		err       error
		errMsg    string
		endpoint  = fmt.Sprintf(`%s/%s`, api.Endpoint, apiMethod)
		payload   = domain3.MessagesSend{
			Message:     `text`,
			AccessToken: cfg.Api.Token,
			ApiVersion:  api.Version,
			PeerId:      peerId,
			RandomId:    rndId,
		}
		testCases = map[string]struct {
			reminder *domain2.Reminder
			keyboard domain3.Keyboard
			withErr  bool
		}{
			`new notifies has 2 buttons: cancel and before5minutes`: {
				reminder: &domain2.Reminder{
					Status: domain2.StatusNew,
				},
				keyboard: domain3.Keyboard{
					OneTime: true,
					Buttons: cancelWith5Minutes(0),
				},
				withErr: false,
			},
			`shifted notifies has only cancel button`: {
				reminder: &domain2.Reminder{
					Status: domain2.StatusShifted,
				},
				keyboard: domain3.Keyboard{
					OneTime: true,
					Buttons: cancel(0),
				},
				withErr: false,
			},
		}
	)

	rnd.On(`Rnd`).Return(rndId)

	for caseName, caseValue := range testCases {
		errMsg = fmt.Sprintf(`test "%s" fall`, caseName)
		kb, _ := json.Marshal(caseValue.keyboard)
		payload.Keyboard = string(kb)
		params, _ := query.Values(payload)
		request, _ = http.NewRequest(`POST`, fmt.Sprintf(`%s?%s`, endpoint, params.Encode()), nil)
		response = &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(answer)),
		}

		client.On(`Do`, request).Return(response, nil)

		err = sm.SendIntention(peerId, payload.Message, caseValue.reminder)
		if (err == nil) == caseValue.withErr {
			t.Errorf(`%s. error result is not expected`, errMsg)
		}
	}
}
