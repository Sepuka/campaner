package analyzer

import (
	"testing"

	"github.com/sepuka/campaner/internal/context"
	"github.com/sepuka/campaner/internal/domain"
	mocks2 "github.com/sepuka/campaner/internal/feature_toggling/toggle/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBadPayload(t *testing.T) {
	var (
		ft     = mocks2.FeatureToggle{}
		logger = zap.NewNop().Sugar()
		msg    = context.Message{
			Payload: ``,
		}
		reminder = &domain.Reminder{}
		err      error
	)
	analyzer := NewAnalyzer(glossary, logger, ft)
	err = analyzer.analyzePayload(msg, reminder)
	assert.NotNil(t, err)
}
