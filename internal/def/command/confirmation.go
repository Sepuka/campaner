package command

import (
	"github.com/sepuka/campaner/internal/command"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const (
	ConfirmationDef = `def.command.confirmation`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: ConfirmationDef,
			Tags: []di.Tag{
				{
					Name: ExecutorDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				return command.NewConfirmation(cfg.Server), nil
			},
		})
	})
}
