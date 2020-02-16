package http

import (
	"net/http"

	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const ClientDef = "http.client"

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: ClientDef,
			Build: func(ctx def.Context) (interface{}, error) {
				return &http.Client{}, nil
			},
		})
	})
}
