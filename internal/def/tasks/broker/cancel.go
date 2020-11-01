package broker

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/repository"
	"github.com/sepuka/campaner/internal/domain"
	"github.com/sepuka/campaner/internal/tasks/broker"
)

const (
	CancelBrokerDef = `def.tasks.broker.cancel`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: CancelBrokerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					storeManager = ctx.Get(repository.ReminderRepoDef).(domain.TaskManager)
				)

				return broker.NewCancelBroker(storeManager), nil
			},
		})
	})
}
