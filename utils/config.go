package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"

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
func (c Config) CheckWithoutLocalPort() error {
	if c.Coin == "ETH" || c.Coin == "ETC" {

		// if c.TCP != 0 {
		// 	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.TCP))
		// 	if err != nil {
		// 		return fmt.Errorf("端口号: %d 转换失败", c.TCP)
		// 	}

		// 	ln, err := net.ListenTCP("tcp", tcpAddr)
		// 	if err != nil {
		// 		return fmt.Errorf("端口号: %d 已经被占用请更换 %s", c.TCP, err.Error())
		// 	}
		// 	defer ln.Close()
		// }

		// if c.TLS != 0 {
		// 	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.TLS))
		// 	if err != nil {
		// 		return fmt.Errorf("端口号: %d 转换失败", c.TLS)
		// 	}

		// 	ln, err := net.ListenTCP("tcp", tcpAddr)
		// 	if err != nil {
		// 		return fmt.Errorf("端口号: %d 已经被占用请更换 %s", c.TLS, err.Error())
		// 	}
		// 	defer ln.Close()
		// }
		//TODO 校验中转矿池是否正确
		if !strings.HasPrefix(c.Pool, "tcp://") && !strings.HasPrefix(c.Pool, "ssl://") {
			return fmt.Errorf("中转矿池地址输入错误(格式为:tcp:// 或 ssl://): %s", c.Pool)
		}

		if c.Mode == 1 {

			return nil
		} else if c.Mode == 2 {
			// if !IsValidHexAddress(c.Wallet) {
			// 	return errors.New("Wallet 钱包地址添加不正确")
			// }
			if !strings.HasPrefix(c.Feepool, "tcp://") && !strings.HasPrefix(c.Feepool, "ssl://") {
				return fmt.Errorf("抽水矿池地址输入错误(格式为:tcp:// 或 ssl://): %s", c.Feepool)
			}
			return nil
		} else {
			return errors.New("不支持的Mode类型")
		}
	} else {
		return errors.New("暂时不支持的币种")
	}
}

// 判断启动参数是否符合要求
func (c Config) Check() error {
	if c.Coin == "ETH" || c.Coin == "ETC" {
		if c.TCP != 0 {
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.TCP))
			if err != nil {
				return fmt.Errorf("端口号: %d 转换失败", c.TCP)
			}

			ln, err := net.ListenTCP("tcp", tcpAddr)
			if err != nil {
				return fmt.Errorf("端口号: %d 已经被占用请更换 %s", c.TCP, err.Error())
			}
			defer ln.Close()
		}

		if c.TLS != 0 {
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.TLS))
			if err != nil {
				return fmt.Errorf("端口号: %d 转换失败", c.TLS)
			}

			ln, err := net.ListenTCP("tcp", tcpAddr)
			if err != nil {
				return fmt.Errorf("端口号: %d 已经被占用请更换 %s", c.TLS, err.Error())
			}
			defer ln.Close()
		}

		// TODO 校验中转矿池是否正确
		if !strings.HasPrefix(c.Pool, "tcp://") && !strings.HasPrefix(c.Pool, "ssl://") {
			return fmt.Errorf("中转矿池地址输入错误(格式为:tcp:// 或 ssl://): %s", c.Pool)
		}

		if c.Mode == 1 {
			return nil
		} else if c.Mode == 2 {
			// if !IsValidHexAddress(c.Wallet) {
			// 	return errors.New("Wallet 钱包地址添加不正确")
			// }

			if !strings.HasPrefix(c.Feepool, "tcp://") && !strings.HasPrefix(c.Feepool, "ssl://") {
				return fmt.Errorf("抽水矿池地址输入错误(格式为:tcp:// 或 ssl://): %s", c.Feepool)
			}
			// if !(strings.HasPrefix("tcp://", c.Feepool) || strings.HasPrefix("ssl://", c.Feepool)) {
			// 	return fmt.Errorf("抽水矿池地址输入错误: 必须已tcp://或ssl:// 开头 %s", c.Feepool)
			// }
			return nil
		} else {
			return errors.New("不支持的Mode类型")
		}
	} else {
		return errors.New("暂时不支持的币种")
	}
}

func (c Config) check_wallet() error {
	return nil
}
