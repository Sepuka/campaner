package config

import (
	"fmt"
	"strings"

	"github.com/stevenroose/gonfig"
)

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

	Database struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     int
	}

	Postpone struct {
		Ids []int
	}

	Features struct {
		Postpone Postpone
	}

	Config struct {
		Server   Server
		Log      Log
		Api      Api
		Db       Database
		Features Features
	}
)

func GetConfig(path string) (*Config, error) {
	var (
		cfg = &Config{}
		err = gonfig.Load(cfg, gonfig.Conf{
			FileDefaultFilename: path,
			FlagIgnoreUnknown:   true,
			FlagDisable:         true,
			EnvDisable:          true,
		})
	)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (api *Api) MaskedToken(params string) string {
	var maskedToken = fmt.Sprintf(`%s...`, api.Token[0:3])

	return strings.Replace(params, api.Token, maskedToken, 1)
}
