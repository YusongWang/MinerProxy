package utils

import (
	"github.com/spf13/viper"
)

const (
	ETH = iota
	ETC
	BTC
	LTC
	CFX
	RVN
	ERGO
	FULX
	ALEO
	FISH
)

type Config struct {
	Coin    string  `mapstructure:"coin"`
	Id      int     `mapstructure:"id"`
	Tcp     int     `mapstructure:"tcp"`
	Tls     int     `mapstructure:"tls"`
	Enport  int     `mapstructure:"enport"`
	Cert    string  `mapstructure:"cert"`
	Key     string  `mapstructure:"key"`
	Pool    string  `mapstructure:"pool"`
	FeePool string  `mapstructure:"feepool"`
	Fee     float64 `mapstructure:"fee"`
	Worker  string  `mapstructure:"worker"`
	Wallet  string  `mapstructure:"wallet"`
	Mode    int     `mapstructure:"mode"`
}

func Parse() Config {
	var config Config
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return config
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config
	}

	return config
}

// 判断启动参数是否符合要求
func (c Config) Check() error {

	return nil
}

func (c Config) check_wallet() error {
	return nil
}
