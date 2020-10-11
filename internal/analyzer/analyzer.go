package analyzer

import (
	"math/rand"

	featureDomain "github.com/sepuka/campaner/internal/feature_toggling/domain"

	"github.com/sepuka/campaner/internal/context"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/domain"
)

type Parser interface {
	Parse(*speeches.Speech, *domain.Reminder) error
	Glossary() []string
	PatternList() []string
}

type Glossary map[string]Parser

type Analyzer struct {
	glossary      Glossary
	logger        *zap.SugaredLogger
	featureToggle featureDomain.FeatureToggle
}

func NewAnalyzer(glossary Glossary, logger *zap.SugaredLogger, feature featureDomain.FeatureToggle) *Analyzer {
	return &Analyzer{
		glossary:      glossary,
		logger:        logger,
		featureToggle: feature,
	}
}

func (a *Analyzer) Analyze(msg context.Message, reminder *domain.Reminder) error {
	var (
		text    = msg.Text
		payload = msg.Payload
	)

	if payload != `` {
		return a.analyzePayload(msg, reminder)
	} else {
		return a.analyzeText(text, reminder)
	}
}

func (a *Analyzer) getRandomStatement(seed int64) string {
	rand.Seed(seed)
	var statements = []string{
		`через 30 минут позвонить другу`,
		`завтра вынести мусор`,
		`вечером сделать домашнюю работу`,
		`в субботу купить корм коту`,
	}

	var rnd = rand.Intn(len(statements) - 1)

	return statements[rnd]
}
