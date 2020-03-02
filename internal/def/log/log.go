package log

import (
	"errors"

	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/def"

	errPkg "github.com/pkg/errors"

	"go.uber.org/zap/zapcore"

	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

const LoggerDef = `logger.def`

func init() {
	def.Register(func(builder *di.Builder, cfg *config.Config) error {
		return builder.Add(di.Def{
			Name: LoggerDef,
			Build: func(ctx def.Context) (interface{}, error) {
				var (
					err               error
					logger            *zap.Logger
					sugar             *zap.SugaredLogger
					zapCfg            zap.Config
					core              zapcore.Core
					fileEncoder       zapcore.Encoder
					fileEncoderConfig zapcore.EncoderConfig
				)

				fileSyncer, closeOut, err := zap.Open(`stdout`)
				if err != nil {
					return nil, errPkg.Wrap(err, `unable to open output files`)
				}

				writeSyncer := zapcore.AddSync(fileSyncer)

				consoleMsgLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					if cfg.Log.Prod {
						return lvl >= zapcore.InfoLevel
					}

					return true
				})

				if cfg.Log.Prod {
					zapCfg = zap.NewProductionConfig()
					fileEncoderConfig = zap.NewProductionEncoderConfig()
				} else {
					zapCfg = zap.NewDevelopmentConfig()
					fileEncoderConfig = zap.NewDevelopmentEncoderConfig()
				}

				zapCfg.OutputPaths = []string{`stdout`}

				fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
				core = zapcore.NewTee(
					zapcore.NewCore(fileEncoder, writeSyncer, consoleMsgLevel),
				)

				logger = zap.New(core)
				sugar = logger.Sugar()
				if sugar == nil {
					closeOut()
					return nil, errors.New(`unable build sugar logger`)
				}

				return sugar, err
			},
			Close: func(obj interface{}) error {
				logger := obj.(*zap.SugaredLogger)
				return logger.Sync()
			},
		})
	})
}
