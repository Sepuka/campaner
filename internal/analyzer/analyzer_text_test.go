package analyzer

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	mocks2 "github.com/sepuka/campaner/internal/feature_toggling/toggle/mocks"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/calendar"

	"github.com/sepuka/campaner/internal/domain"

	"github.com/stretchr/testify/assert"
)

var (
	parsers = []Parser{
		NewTimeParser(),
		NewDayParser(),
		NewDateParser(),
		NewDateMonthParser(),
		NewTimesOfDayParser(),
		NewDateTimeAggregateParser([]Parser{NewTimeParser(), NewDayParser()}),
	}
	glossary = make(map[string]Parser)
)

func buildGlossary() {
	var (
		keyword string
	)
	for _, parser := range parsers {
		for _, keyword = range parser.(Parser).Glossary() {
			glossary[keyword] = parser.(Parser)
		}
	}
}
func TestMain(m *testing.M) {
	buildGlossary()
	code := m.Run()
	os.Exit(code)
}

func TestNewAnalyzer(t *testing.T) {
	var (
		ft     = mocks2.FeatureToggle{}
		logger = zap.NewNop().Sugar()
		err    error
	)
	analyzer := NewAnalyzer(glossary, logger, ft)

	var testCases = map[string]struct {
		speech           string
		expectedReminder *domain.Reminder
	}{
		`напомни мне через 25 секунд что-то сделать`: {
			speech: `напомни мне через 25 секунд что-то сделать`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`напомни мне что-то сделать`, ` `),
				When:    time.Duration(25) * time.Second,
			},
		},
		`напомни в 23:15 что-то сделать`: {
			speech: `напомни В 23:15 что-то сделать`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`напомни что-то сделать`, ` `),
				When:    time.Until(calendar.NextNight().Add(15 * time.Minute)),
			},
		},
		`завтра в 09:23 отвести детей в школу`: {
			speech: `завтра в 09:23 отвести детей в школу`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`отвести детей в школу`, ` `),
				When:    time.Until(calendar.NextMidnight().Add(9 * time.Hour).Add(23 * time.Minute)),
			},
		},
		`через час хлеба и зрелищ`: {
			speech: `через час хлеба и зрелищ`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`хлеба и зрелищ`, ` `),
				When:    time.Hour,
			},
		},
		`сегодня починить кофеварку`: {
			speech: `сегодня починить кофеварку`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`починить кофеварку`, ` `),
				When:    calendar.GetNextPeriod(calendar.NewDate(time.Now())).Until(),
			},
		},
		`утром`: {
			speech: `утром купить хлеба`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`купить хлеба`, ` `),
				When:    time.Until(calendar.NextMorning()),
			},
		},
		`днем`: {
			speech: `днем`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`днем`, ` `),
				When:    time.Until(calendar.NextAfternoon()),
			},
		},
		`вечером`: {
			speech: `вечером`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`вечером`, ` `),
				When:    time.Until(calendar.NextEvening()),
			},
		},
		`ночью`: {
			speech: `ночью`,
			expectedReminder: &domain.Reminder{
				Subject: strings.Split(`ночью`, ` `),
				When:    time.Until(calendar.NextNight()),
			},
		},
	}

	for testName, testCase := range testCases {
		var (
			testError        = fmt.Sprintf(`test "%s" error`, testName)
			expectedReminder = testCase.expectedReminder
			actualReminder   = domain.NewReminder(0)
		)
		err = analyzer.analyzeText(testCase.speech, actualReminder)
		assert.NoError(t, err)
		assert.InDelta(t, expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.Equal(t, expectedReminder.GetSubject(), actualReminder.GetSubject(), testError)
	}
}

func TestDayOfWeekAnalyzer(t *testing.T) {
	var (
		ft     = mocks2.FeatureToggle{}
		logger = zap.NewNop().Sugar()
	)
	analyzer := NewAnalyzer(glossary, logger, ft)
	var testCases = map[string]struct {
		words    string
		reminder *domain.Reminder
	}{
		`в понедельник и время с минутами`: {
			words: `в понедельник в 16:00 часов встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextMonday().Add(time.Hour * 16).Until(),
			},
		},
		`во вторник и время с минутами`: {
			words: `во вторник в 17:00 часов встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextTuesday().Add(time.Hour * 17).Until(),
			},
		},
		`в среду и время с минутами`: {
			words: `в среду в 18:00 часов встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextWednesday().Add(time.Hour * 18).Until(),
			},
		},
		`в четверг и время с минутами`: {
			words: `в четверг в 19:00 часов встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextThursday().Add(time.Hour * 19).Until(),
			},
		},
		`в пятницу и время с минутами`: {
			words: `в пятницу в 20:00 часов встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextFriday().Add(time.Hour * 20).Until(),
			},
		},
		`в субботу и время с минутами`: {
			words: `в субботу в 21:00 час встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextSaturday().Add(time.Hour * 21).Until(),
			},
		},
		`в воскресенье и время с минутами`: {
			words: `в воскресенье в 22:00 часа встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`встреча`},
				When:    calendar.NextSunday().Add(time.Hour * 22).Until(),
			},
		},
		`день недели и время без минут`: {
			words: `В среду в 16 пройдет маленькая пятничная встреча`,
			reminder: &domain.Reminder{
				Subject: []string{`пройдет`, `маленькая`, `пятничная`, `встреча`},
				When:    calendar.NextWednesday().Add(time.Hour * 16).Until(),
			},
		},
	}

	for testName, testCase := range testCases {
		var (
			testError        = fmt.Sprintf(`test "%s" error`, testName)
			expectedReminder = testCase.reminder
			actualReminder   = &domain.Reminder{}
		)
		analyzer.analyzeText(testCase.words, actualReminder)
		assert.InDelta(t, expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.Equal(t, expectedReminder.GetSubject(), actualReminder.GetSubject(), testError)
	}
}

func TestDateAnalyzer(t *testing.T) {
	var (
		logger         = zap.NewNop().Sugar()
		ft             = mocks2.FeatureToggle{}
		analyzer       = NewAnalyzer(glossary, logger, ft)
		now            = time.Now()
		futureMidnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(calendar.Day * 3)
		futureMoment   = fmt.Sprintf(`%s в 18:00 собрание`, futureMidnight.Format(calendar.DayMonthFormat))
	)

	var testCases = map[string]struct {
		words    string
		reminder *domain.Reminder
	}{
		`указано время и дата`: {
			words: futureMoment,
			reminder: &domain.Reminder{
				Subject: []string{`собрание`},
				When:    time.Until(futureMidnight.Add(18 * time.Hour)),
			},
		},
	}

	for testName, testCase := range testCases {
		var (
			testError        = fmt.Sprintf(`test "%s" error`, testName)
			expectedReminder = testCase.reminder
			actualReminder   = &domain.Reminder{}
		)
		analyzer.analyzeText(testCase.words, actualReminder)
		assert.InDelta(t, expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.Equal(t, expectedReminder.GetSubject(), actualReminder.GetSubject(), testError)
	}
}

func TestDateMonthAnalyzer(t *testing.T) {
	var (
		now               = time.Now()
		lastDayOfThisYear = time.Date(now.Year(), 12, 31, 9, 0, 0, 0, time.Local)
		actualErr         error
	)

	var testCases = map[string]struct {
		words    string
		reminder *domain.Reminder
		wantErr  bool
	}{
		`31 декабря корпоративная вечеринка`: {
			words: `31 декабря корпоративная вечеринка`,
			reminder: &domain.Reminder{
				Subject: strings.Split(`корпоративная вечеринка`, ` `),
				When:    time.Until(lastDayOfThisYear),
			},
			wantErr: false,
		},
		`дата дублирована`: {
			words: `31 31 декабря корпоративная вечеринка`,
			reminder: &domain.Reminder{
				Subject: strings.Split(`корпоративная вечеринка`, ` `),
				When:    time.Until(lastDayOfThisYear),
			},
			wantErr: false,
		},
	}

	for testName, testCase := range testCases {
		var (
			testError        = fmt.Sprintf(`test "%s" error`, testName)
			expectedReminder = testCase.reminder
			actualReminder   = &domain.Reminder{}
			logger           = zap.NewNop().Sugar()
			ft               = mocks2.FeatureToggle{}
			analyzer         = NewAnalyzer(glossary, logger, ft)
		)

		actualErr = analyzer.analyzeText(testCase.words, actualReminder)
		if testCase.wantErr && actualErr == nil {
			t.Errorf("Parse() error = %v, wantErr %v", actualErr, testCase.wantErr)
			return
		}
		assert.InDelta(t, expectedReminder.When.Seconds(), actualReminder.When.Seconds(), 1, testError)
		assert.Equal(t, expectedReminder.GetSubject(), actualReminder.GetSubject(), testError)
	}
}

func TestUnknownTime(t *testing.T) {
	var (
		logger      = zap.NewNop().Sugar()
		ft          = mocks2.FeatureToggle{}
		analyzer    = NewAnalyzer(glossary, logger, ft)
		reminder    = &domain.Reminder{}
		actualError error
	)

	actualError = analyzer.analyzeText(`abc def`, reminder)
	assert.Error(t, actualError)
	assert.Equal(t, float64(0), reminder.When.Seconds())
}
