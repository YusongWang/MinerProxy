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
	Config      []utils.Config `json:"config"`
	Pools       [][]pack.Worker
	Port        int
	Password    string
	Jwt_secret  string
}

var WebApp = new(WebApplication)

type ManageConfig struct {
	Config []utils.Config `json:"config"`
	Web    struct {
		Port     int    `json:"port"`
		Password string `json:"password"`
	} `json:"web"`
}

var ManageApp = new(ManageConfig)

var OnlinePools [1000]map[string]pack.Worker
