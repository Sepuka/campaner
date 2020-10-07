package analyzer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	featureDomain "github.com/sepuka/campaner/internal/feature_toggling/domain"

	"github.com/sepuka/campaner/internal/context"

	domain2 "github.com/sepuka/campaner/internal/api/domain"
	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/sepuka/campaner/internal/errors"

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
	taskManager   domain.TaskManager
	featureToggle featureDomain.FeatureToggle
}

func NewAnalyzer(glossary Glossary, logger *zap.SugaredLogger, taskManager domain.TaskManager, feature featureDomain.FeatureToggle) *Analyzer {
	return &Analyzer{
		glossary:      glossary,
		logger:        logger,
		taskManager:   taskManager,
		featureToggle: feature,
	}
}

func (a *Analyzer) Analyze(msg context.Message, reminder *domain.Reminder) {
	var (
		text    = msg.Text
		payload = msg.Payload
	)

	if payload != `` {
		a.analyzePayload(msg, reminder)
	} else {
		a.analyzeText(text, reminder)
	}
}

func (a *Analyzer) analyzeText(text string, reminder *domain.Reminder) {
	a.buildReminder(speeches.NewSpeech(text), reminder)
}

func (a *Analyzer) analyzePayload(msg context.Message, reminder *domain.Reminder) {
	var (
		payload    domain2.ButtonPayload
		err        error
		taskId     int64
		rawPayload = msg.Payload
		text       = domain2.ButtonText(msg.Text)
	)

	if err = json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		a.logger.
			With(
				zap.String(`payload`, rawPayload),
				zap.Error(err),
			).
			Error(`analyze payload error`)
		return
	}

	if taskId, err = strconv.ParseInt(payload.Button, 10, 64); err != nil {
		a.
			logger.
			With(
				zap.String(`json`, rawPayload),
				zap.Int(`user_id`, reminder.Whom),
				zap.Error(err),
			).
			Error(`cannot parse task_id`)
		return
	}

	switch text {
	case domain2.CancelButton:
		if err = a.taskManager.Cancel(taskId, reminder.Whom); err != nil {
			a.
				logger.
				With(
					zap.Int64(`task_id`, taskId),
					zap.Int(`user_id`, reminder.Whom),
					zap.Error(err),
				).
				Error(`cannot cancel task`)
			return
		}
		reminder.Subject = []string{`напоминание отменено`}
		reminder.When = time.Nanosecond
	case domain2.Later15MinButton:
		if !a.featureToggle.IsEnabled(reminder.Whom, featureDomain.Postpone) {
			return
		}
		var minutes int
		switch text {
		case domain2.Later15MinButton:
			minutes = 15
		default:
			minutes = 30
		}
		if err = a.taskManager.Prolong(taskId, reminder.Whom, minutes); err != nil {
			a.
				logger.
				With(
					zap.Int64(`task_id`, taskId),
					zap.Int(`user_id`, reminder.Whom),
					zap.Error(err),
				).
				Error(`cannot prolong task`)
			return
		}
		reminder.Subject = []string{`напоминание продлено`}
		reminder.When = time.Nanosecond
	}
}

func (a *Analyzer) buildReminder(speech *speeches.Speech, reminder *domain.Reminder) {
	const patternLength = 1
	var (
		err     error
		pattern *speeches.Pattern
	)

	if pattern, err = speech.TryPattern(patternLength); err != nil {
		if reminder.GetSubject() == `` {
			reminder.AppendSubject(speeches.NewPattern([]string{`ring!`}))
		}
		if reminder.IsTimeUnknown() {
			var randomSubject = fmt.Sprintf(`Попробуйте фразу: "%s"`, a.getRandomStatement(time.Now().Unix()))
			reminder.RewriteSubject(randomSubject)
			reminder.When = time.Second
		}
		return
	}

	if parser, ok := a.glossary[pattern.Origin()]; ok {
		if err = parser.Parse(speech, reminder); err != nil {
			var (
				patterns, what string
			)

			switch errors.GetType(err) {
			case errors.ItIsPastTimeError:
				what = `it is past time!`
			default:
				patterns = strings.Join(parser.PatternList(), "\n")
				what = fmt.Sprintf("use known format, for instance:\n%s\n", patterns)
			}

			*reminder = *domain.NewImmediateReminder(reminder.Whom, what)
			return
		}
	} else {
		if err = speech.ApplyPattern(pattern); err != nil {
			return
		}
		reminder.AppendSubject(pattern)
	}

	a.buildReminder(speech, reminder)
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
