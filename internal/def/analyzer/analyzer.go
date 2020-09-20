package analyzer

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/log"
	"github.com/sepuka/campaner/internal/def/repository"
	"github.com/sepuka/campaner/internal/domain"
	"go.uber.org/zap"
)

const (
	AnalyzerDef  = `def.analyzer`
	ParserTagDef = `dev.analyzer.parser`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: AnalyzerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					parsers  = def.GetByTag(ParserTagDef)
					glossary = make(analyzer.Glossary)
					keyword  string
					logger   = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					repo     = ctx.Get(repository.ReminderRepoDef).(domain.TaskManager)
				)

				for _, parser := range parsers {
					for _, keyword = range parser.(analyzer.Parser).Glossary() {
						glossary[keyword] = parser.(analyzer.Parser)
					}
				}

				return analyzer.NewAnalyzer(glossary, logger, repo), nil
			},
		})
	})
}
