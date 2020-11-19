package broker

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/tasks/broker"
)

const (
	BarrenBrokerDef = `def.tasks.broker.barren`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: BarrenBrokerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				return broker.NewBarrenBroker(), nil
			},
		})
	})
}
