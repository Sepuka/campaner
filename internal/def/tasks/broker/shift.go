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
	ShiftBrokerDef = `def.tasks.broker.shift`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: ShiftBrokerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					taskManager = ctx.Get(repository.ReminderRepoDef).(domain.TaskManager)
					repo        = ctx.Get(repository.ReminderRepoDef).(domain.ReminderRepository)
				)

				return broker.NewShiftBroker(taskManager, repo), nil
			},
		})
	})
}
