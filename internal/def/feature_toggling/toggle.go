package feature_toggling

import (
	"github.com/sepuka/campaner/internal/feature_toggling/toggle"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const FeatureToggleDef = "feature.toggler"

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: FeatureToggleDef,
			Build: func(ctx def.Context) (interface{}, error) {

				return toggle.NewToggle(cfg), nil
			},
		})
	})
}
