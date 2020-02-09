package bot

import (
	command2 "github.com/sepuka/campaner/internal/def/command"
	"github.com/sepuka/campaner/internal/def/log"
	"github.com/sepuka/campaner/internal/def/middleware"
	middleware2 "github.com/sepuka/campaner/internal/middleware"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/bot"
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"go.uber.org/zap"
)

const (
	BotDef = `def.bot`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: BotDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					logger     = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					handlerMap = ctx.Get(command2.HandlerMapDef).(command.HandlerMap)
					handler    = ctx.Get(middleware.BotMiddlewareDef).(middleware2.HandlerFunc)
				)

				return bot.NewBot(logger, handlerMap, handler, cfg.Server), nil
			},
		})
	})
}
