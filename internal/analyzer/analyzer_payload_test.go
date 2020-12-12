package analyzer

import (
	"fmt"
	"testing"

	"github.com/sepuka/campaner/internal/api/method"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	mocks2 "github.com/sepuka/campaner/internal/feature_toggling/toggle/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPayloadWithTaskButton(t *testing.T) {
	var (
		ft       = mocks2.FeatureToggle{}
		logger   = zap.NewNop().Sugar()
		reminder *domain.Reminder
		err      error

		testCases = map[string]struct {
			msg    context.Message
			status int
			err    bool
		}{
			`normal cancel`: {
				msg: context.Message{
					Text:    `there is can be any text, but "cancel" is prefer`,
					Payload: `{"button":"12345", "command":"cancel"}`,
				},
				status: domain.StatusCanceled,
				err:    false,
			},
			`cancel without taskId`: {
				msg: context.Message{
					Text:    `there is can be any text, but "cancel" is prefer`,
					Payload: `{"command":"cancel"}`,
				},
				status: domain.StatusNew,
				err:    true,
			},
			`normal 15 minutes`: {
				msg: context.Message{
					Text:    `there is can be any text, but "later 15 minutes" is prefer`,
					Payload: `{"button":"12345", "command":"after15minutes"}`,
				},
				status: domain.StatusCopied,
				err:    false,
			},
			`15 minutes without taskId`: {
				msg: context.Message{
					Text:    string(method.Later15MinButton),
					Payload: ``,
				},
				status: domain.StatusNew,
				err:    true,
			},
			`normal OK button`: {
				msg: context.Message{
					Text:    `there is can be any text, but OK is prefer`,
					Payload: `{"button":"12345", "command":"OK"}`,
				},
				status: domain.StatusBarren,
				err:    false,
			},
		}
	)

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		reminder = &domain.Reminder{}

		analyzer := NewAnalyzer(glossary, logger, ft)
		err = analyzer.analyzePayload(testCase.msg, reminder)
		if (err != nil) != testCase.err {
			t.Errorf(`unexpected error %v in %s`, err, testError)
			return
		}

		assert.Equal(t, testCase.status, reminder.Status)
	}
}

func TestPayloadWithStartButton(t *testing.T) {
	var (
		ft       = mocks2.FeatureToggle{}
		logger   = zap.NewNop().Sugar()
		reminder *domain.Reminder
		err      error

		testCases = map[string]struct {
			msg    context.Message
			status int
			err    bool
		}{
			`normal start`: {
				msg: context.Message{
					Text:    `there is can be any text, but "Начать" is prefer`,
					Payload: `{"command":"start"}`,
				},
				status: domain.StatusBarren,
				err:    false,
			},
			`normal start and some button ID (we don'\t mind about button ID)`: {
				msg: context.Message{
					Text:    `there is can be any text, but "Начать" is prefer`,
					Payload: `{"command":"start", "button": "111"}`,
				},
				status: domain.StatusBarren,
				err:    false,
			},
		}
	)

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		reminder = &domain.Reminder{}

		analyzer := NewAnalyzer(glossary, logger, ft)
		err = analyzer.analyzePayload(testCase.msg, reminder)
		if (err != nil) != testCase.err {
			t.Errorf(`unexpected error %v in %s`, err, testError)
			return
		}

		assert.Equal(t, testCase.status, reminder.Status)
	}
}
