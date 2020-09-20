package command

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const (
	HandlerMapDef = `handler.map.def`
	ExecutorDef   = `handler.command.def`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: HandlerMapDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					handlers   = def.GetByTag(ExecutorDef)
					handlerMap = make(command.HandlerMap, len(handlers))
					precept    string
				)

				for _, cmd := range handlers {
					for _, precept = range cmd.(command.Preceptable).Precept() {
						handlerMap[precept] = cmd.(command.Executor)
					}
				}

				return handlerMap, nil
			},
		})
	})
}
