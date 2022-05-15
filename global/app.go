package global

import (
	"miner_proxy/utils"
	"time"

	"github.com/spf13/viper"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      []utils.Config
	Pools       [][]Worker
}

var App = new(Application)

type WebApplication struct {
	ConfigViper *viper.Viper
	Config      []utils.Config `json:"config"`
	Pools       [][]Worker
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

var OnlinePools [1000]map[string]Worker

var Start_Time = time.Now().Unix()

func init() {
	Start_Time = time.Now().Unix()
}
