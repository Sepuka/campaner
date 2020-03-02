package bot

import (
	"github.com/sarulabs/di"
	api3 "github.com/sepuka/campaner/internal/api"
	"github.com/sepuka/campaner/internal/bot"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	api2 "github.com/sepuka/campaner/internal/def/api"
	"github.com/sepuka/campaner/internal/def/log"
	repository2 "github.com/sepuka/campaner/internal/def/repository"
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
					repo   = ctx.Get(repository2.ReminderRepoDef).(domain.ReminderRepository)
					logger = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					api    = ctx.Get(api2.SendMessageDef).(*api3.SendMessage)
				)

				return bot.NewWorker(repo, logger, api), nil
			},
		})
	})
}
