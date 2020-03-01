package command

import (
	"github.com/sarulabs/di"
	analyzer3 "github.com/sepuka/campaner/internal/analyzer"
	api3 "github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	analyzer2 "github.com/sepuka/campaner/internal/def/analyzer"
	api2 "github.com/sepuka/campaner/internal/def/api"
	"github.com/sepuka/campaner/internal/def/log"
	"github.com/sepuka/campaner/internal/def/repository"
	"github.com/sepuka/campaner/internal/domain"
	"go.uber.org/zap"
)

const (
	MessageNewDef = `def.command.message_new`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: MessageNewDef,
			Tags: []di.Tag{
				{
					Name: ExecutorDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					api      = ctx.Get(api2.SendMessageDef).(*api3.SendMessage)
					logger   = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					analyzer = ctx.Get(analyzer2.AnalyzerDef).(*analyzer3.Analyzer)
					repo     = ctx.Get(repository.ReminderRepoDef).(domain.ReminderRepository)
				)

				return command.NewMessageNew(cfg.Server, api, logger, analyzer, repo), nil
			},
		})
	})
}
