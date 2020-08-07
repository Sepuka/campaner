package api

import (
	http2 "net/http"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/api/method"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/http"
	"github.com/sepuka/campaner/internal/def/log"
	"go.uber.org/zap"
)

const (
	UsersGetDef = `def.api.users_get`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: UsersGetDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					client = ctx.Get(http.ClientDef).(*http2.Client)
					logger = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
				)
				return method.NewUsersGet(cfg, client, logger), nil
			},
		})
	})
}
