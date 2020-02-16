package config

import "github.com/stevenroose/gonfig"

type (
	Server struct {
		Server       string
		Confirmation string
		Socket       string `default:"/var/run/campaner.sock"`
	}

	Log struct {
		Prod   bool
		Output string
	}

	Api struct {
		Token string
	}

	Config struct {
		Server Server
		Log    Log
		Api    Api
	}
)

func GetConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := gonfig.Load(cfg, gonfig.Conf{
		FileDefaultFilename: path,
		FlagIgnoreUnknown:   true,
		FlagDisable:         true,
		EnvDisable:          true,
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
