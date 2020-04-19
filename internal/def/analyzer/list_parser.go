package analyzer

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/repository"
	"github.com/sepuka/campaner/internal/domain"
)

const (
	ListParserDef = `def.analyzer.parser.list`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: ListParserDef,
			Tags: []di.Tag{
				{
					Name: ParserTagDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				repo := ctx.Get(repository.ReminderRepoDef).(domain.ReminderRepository)

				return analyzer.NewListParser(repo), nil
			},
		})
	})
}
