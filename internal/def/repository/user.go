package repository

import (
	"github.com/go-pg/pg"
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	db2 "github.com/sepuka/campaner/internal/def/db"
	"github.com/sepuka/campaner/internal/repository"
)

const UserRepoDef = `repo.user.def`

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: UserRepoDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					db = ctx.Get(db2.DataBaseDef).(*pg.DB)
				)

				return repository.NewUserRepository(db), nil
			},
		})
	})
}
