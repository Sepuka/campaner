package analyzer

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const (
	DayParserDef = `def.analyzer.parser.day`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: DayParserDef,
			Tags: []di.Tag{
				{
					Name: ParserTagDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				return analyzer.NewDayParser(), nil
			},
		})
	})
}
