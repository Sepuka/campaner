package analyzer

import (
	"github.com/sarulabs/di"
	"github.com/sepuka/campaner/internal/analyzer"
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"
)

const (
	DateTimeAggregateParserDef = `def.analyzer.parser.date_time_aggregate`
)

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: DateTimeAggregateParserDef,
			Tags: []di.Tag{
				{
					Name: ParserTagDef,
					Args: nil,
				},
			},
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					timeParser   = ctx.Get(TimeParserDef).(analyzer.Parser)
					dayParser    = ctx.Get(DayParserDef).(analyzer.Parser)
					periodParser = ctx.Get(PeriodParserDef).(analyzer.Parser)
					parsers      = []analyzer.Parser{
						timeParser,
						dayParser,
						periodParser,
					}
				)

				return analyzer.NewDateTimeAggregateParser(parsers), nil
			},
		})
	})
}
