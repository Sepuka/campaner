package tasks

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/tasks/broker"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/tasks"
)

const (
	BrokerDef = `def.tasks.broker`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: BrokerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					brokerMap = tasks.BrokerMap{
						domain.StatusNew:      ctx.Get(broker.NewBornDef).(tasks.Broker),
						domain.StatusCopied:   ctx.Get(broker.CopyBrokerDef).(tasks.Broker),
						domain.StatusCanceled: ctx.Get(broker.CancelBrokerDef).(tasks.Broker),
						domain.StatusShifted:  ctx.Get(broker.ShiftBrokerDef).(tasks.Broker),
						domain.StatusBarren:   ctx.Get(broker.BarrenBrokerDef).(tasks.Broker),
					}
				)

				return tasks.NewTaskBroker(brokerMap), nil
			},
		})
	})
}
