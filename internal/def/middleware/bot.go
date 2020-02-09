package middleware

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/middleware"
)

const (
	BotMiddlewareDef = `middleware.bot.def`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: BotMiddlewareDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					terminalMiddleware = []func(handlerFunc middleware.HandlerFunc) middleware.HandlerFunc{
						middleware.Panic,
					}
				)

				return middleware.BuildHandlerChain(terminalMiddleware), nil
			},
		})
	})
}
