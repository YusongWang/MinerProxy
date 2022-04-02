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
	Coin    string
	TcpPort int
	SslPort int
	EnPort  int
	Cert    string
	Key     string
	Pool    string
	FeePool string
	Fee     float32
}

func Parse() error {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		//viper.Read
		viper.SetEnvPrefix("MinerProxy_")
	}

	return nil
}
