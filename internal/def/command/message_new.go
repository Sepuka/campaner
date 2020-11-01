package command

import (
	"github.com/sarulabs/di"
	analyzer3 "github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/api/method"
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	analyzer2 "github.com/sepuka/campaner/internal/def/analyzer"
	api2 "github.com/sepuka/campaner/internal/def/api"
	"github.com/sepuka/campaner/internal/def/log"
	tasks2 "github.com/sepuka/campaner/internal/def/tasks"
	"github.com/sepuka/campaner/internal/tasks"
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
					api          = ctx.Get(api2.SendMessageDef).(*method.SendMessage)
					logger       = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					analyzer     = ctx.Get(analyzer2.AnalyzerDef).(*analyzer3.Analyzer)
					storeManager = ctx.Get(tasks2.BrokerDef).(*tasks.TaskBroker)
				)

				return command.NewMessageNew(api, logger, analyzer, storeManager), nil
			},
		})
	})
}
