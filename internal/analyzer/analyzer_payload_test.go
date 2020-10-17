package analyzer

import (
	"fmt"
	"testing"

	domain2 "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	mocks2 "github.com/sepuka/campaner/internal/feature_toggling/toggle/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPayloadWithCancelButton(t *testing.T) {
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
					Text:    string(domain2.CancelButton),
					Payload: `{"button":"12345"}`,
				},
				status: domain.StatusCanceled,
				err:    false,
			},
			`cancel without taskId`: {
				msg: context.Message{
					Text:    string(domain2.CancelButton),
					Payload: ``,
				},
				status: domain.StatusNew,
				err:    true,
			},
			`normal 15 minutes`: {
				msg: context.Message{
					Text:    string(domain2.Later15MinButton),
					Payload: `{"button":"12345"}`,
				},
				status: domain.StatusCopied,
				err:    false,
			},
			`15 minutes without taskId`: {
				msg: context.Message{
					Text:    string(domain2.Later15MinButton),
					Payload: ``,
				},
				status: domain.StatusNew,
				err:    true,
			},
			`normal OK button`: {
				msg: context.Message{
					Text:    string(domain2.OKButton),
					Payload: `{"button":"12345"}`,
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
