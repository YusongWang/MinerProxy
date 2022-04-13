package global

import (
	"miner_proxy/pack"
	"miner_proxy/utils"

	"github.com/spf13/viper"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      []utils.Config
	Pools       [][]pack.Worker
}

var App = new(Application)

type WebApplication struct {
	ConfigViper *viper.Viper
	Config      []utils.Config
	Pools       [][]pack.Worker
}

var WebApp = new(WebApplication)
