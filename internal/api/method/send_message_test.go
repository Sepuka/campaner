package method

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sepuka/campaner/internal/context"

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

	var sm = NewSendMessage(cfg, client, logger, ft)
	err = sm.SendNotification(peerId, ``, 1)
	assert.NoError(t, err)
}
