package config

import "github.com/stevenroose/gonfig"

type (
	Server struct {
		Server       string
		Confirmation string
		Socket       string `default:"/var/run/campaner.sock"`
	}

	Log struct {
		Prod bool
	}

	Api struct {
		Token string `default:"???_there_is_the_access_api_token"`
	}

	Config struct {
		Server Server
		Log    Log
		Api    Api
		Db     Database
	}

	Database struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     int
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
