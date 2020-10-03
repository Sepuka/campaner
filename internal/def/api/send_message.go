package api

import (
	http2 "net/http"

	"github.com/sepuka/campaner/internal/def/feature_toggling"
	featureDomain "github.com/sepuka/campaner/internal/feature_toggling/domain"

	"github.com/sepuka/campaner/internal/api/method"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/http"
	"github.com/sepuka/campaner/internal/def/log"
	"go.uber.org/zap"
)

const (
	SendMessageDef = `def.api.send_message`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: SendMessageDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					client        = ctx.Get(http.ClientDef).(*http2.Client)
					logger        = ctx.Get(log.LoggerDef).(*zap.SugaredLogger)
					featureToggle = ctx.Get(feature_toggling.FeatureToggleDef).(featureDomain.FeatureToggle)
				)
				return method.NewSendMessage(cfg, client, logger, featureToggle), nil
			},
		})
	})
}
