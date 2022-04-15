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
	Coin    string  `json:"coin"`
	ID      int     `json:"id"`
	TCP     int     `json:"tcp"`
	TLS     int     `json:"tls"`
	Enport  int     `json:"enport"`
	Cert    string  `json:"cert"`
	Key     string  `json:"key"`
	Pool    string  `json:"pool"`
	Feepool string  `json:"feepool"`
	Fee     float64 `json:"fee"`
	Worker  string  `json:"worker"`
	Wallet  string  `json:"wallet"`
	Mode    int     `json:"mode"`
	Online  bool    `json:"online"`
	// Coin    string  `json:"Coin"`
	// ID      int     `json:"ID"`
	// TCP     int     `json:"TCP"`
	// TLS     int     `json:"TLS"`
	// Enport  int     `json:"Enport"`
	// Cert    string  `json:"Cert"`
	// Key     string  `json:"Key"`
	// Pool    string  `json:"Pool"`
	// Feepool string  `json:"Feepool"`
	// Fee     float64 `json:"Fee"`
	// Worker  string  `json:"Worker"`
	// Wallet  string  `json:"Wallet"`
	// Mode    int     `json:"Mode"`
}

func Parse() Config {
	var config Config
	viper.SetConfigName("config.yaml")    // name of config file (without extension)
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
