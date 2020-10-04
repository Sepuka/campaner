package bot

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/api/method"
	"github.com/sepuka/campaner/internal/bot"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/api"
	"github.com/sepuka/campaner/internal/def/log"
	"github.com/sepuka/campaner/internal/def/repository"
	"github.com/sepuka/campaner/internal/domain"
	"go.uber.org/zap"
)

const (
	WorkerDef = `def.bot.worker`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: WorkerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					repo      = ctx.Get(repository.ReminderRepoDef).(domain.ReminderRepository)
					logger    = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					messenger = ctx.Get(api.SendMessageDef).(*method.SendMessage)
				)

				return bot.NewWorker(repo, logger, messenger), nil
			},
		})
	})
}
