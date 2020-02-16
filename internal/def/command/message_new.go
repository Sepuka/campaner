package command

import (
	"github.com/sarulabs/di"
	api3 "github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	api2 "github.com/sepuka/campaner/internal/def/api"
	"github.com/sepuka/campaner/internal/def/log"
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
					api    = ctx.Get(api2.SendMessageDef).(*api3.SendMessage)
					logger = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
				)

				return command.NewMessageNew(cfg.Server, api, logger), nil
			},
		})
	})
}
