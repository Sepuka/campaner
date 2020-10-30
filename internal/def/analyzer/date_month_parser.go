package analyzer

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const (
	DateMonthParserDef = `def.analyzer.parser.date_month`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: DateMonthParserDef,
			Tags: []di.Tag{
				{
					Name: ParserTagDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				return analyzer.NewDateMonthParser(), nil
			},
		})
	})
}
