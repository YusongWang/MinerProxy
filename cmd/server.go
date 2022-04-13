package cmd

import (
	//_ "net/http/pprof"
	"os"

	etcboot "miner_proxy/boot/etc"
	ethboot "miner_proxy/boot/eth"
	"miner_proxy/boot/eth_test"
	"miner_proxy/boot/test"
	"miner_proxy/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	viper.AutomaticEnv()

	serverCmd.Flags().String("coin", "ETH", "指定需要代理的币种")
	viper.BindPFlag("coin", serverCmd.Flags().Lookup("coin"))

	serverCmd.Flags().Int("id", 0, "指定当前代理编号")
	viper.BindPFlag("id", serverCmd.Flags().Lookup("id"))

	serverCmd.Flags().String("crt", "cert.pem", "指定SSL服务器证书")
	viper.BindPFlag("crt", serverCmd.Flags().Lookup("crt"))

	serverCmd.Flags().String("key", "key.pem", "指定SSL证书私钥")
	viper.BindPFlag("key", serverCmd.Flags().Lookup("key"))

	serverCmd.Flags().Int("tcp", 0, "指定需要代理的TCP端口号")
	viper.BindPFlag("tcp", serverCmd.Flags().Lookup("tcp"))

	serverCmd.Flags().Int("tls", 0, "指定需要代理的TLS端口号")
	viper.BindPFlag("tls", serverCmd.Flags().Lookup("tls"))

	serverCmd.Flags().Int("encrypt", 0, "指定需要代理的加密服务端口号")
	viper.BindPFlag("encrypt", serverCmd.Flags().Lookup("encrypt"))

	serverCmd.Flags().String("wallet", "", "抽水钱包地址")
	viper.BindPFlag("wallet", serverCmd.Flags().Lookup("wallet"))

	serverCmd.Flags().String("worker", "MinerProxy", "抽水矿工名称")
	viper.BindPFlag("worker", serverCmd.Flags().Lookup("worker"))

	serverCmd.Flags().String("pool", "", "指定需要配置代理中转的矿池地址\n格式: (ssl://)tcp://asia2.ethermine.org:4444")
	viper.BindPFlag("pool", serverCmd.Flags().Lookup("pool"))

	serverCmd.Flags().String("feepool", "", "指定抽水的矿池地址\n格式： (ssl://)tcp://asia2.ethermine.org:4444")
	viper.BindPFlag("feepool", serverCmd.Flags().Lookup("feepool"))

	serverCmd.Flags().Int("mode", 0, "中转模式: 1. 直连 2. 抽水")
	viper.BindPFlag("mode", serverCmd.Flags().Lookup("mode"))

	serverCmd.Flags().Float64("fee", 0.0, "抽水率(%)默认 2% 支持一位小数点。")
	viper.BindPFlag("fee", serverCmd.Flags().Lookup("fee"))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "开启纯转发模式，不启动web界面。",
	Long:  `开启纯转发模式，不启动web界面。`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(pflag.CommandLine)
		config := parseConfig()

		if err := config.Check(); err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}

		// go func() {
		// 	log.Println(http.ListenAndServe(":6060", nil))
		// }()

		if config.Mode == 1 {
			switch config.Coin {
			case "ETH":
				ethboot.BootNoFee(config)
			case "ETC":
				etcboot.BootNoFee(config)
			default:
				test.BootNoFee(config)
			}
		} else if config.Mode == 2 {
			switch config.Coin {
			case "ETH":
				ethboot.BootWithFee(config)
			case "ETC":
				etcboot.BootWithFee(config)
			case "ETH_TEST":
				eth_test.BootWithFee(config)
			default:
				test.BootNoFee(config)
			}
		} else {
			utils.Logger.Error("不支持的Mode参数")
			os.Exit(99)
		}
	},
}

func parseFromCli(c *utils.Config) {
	viper.SetEnvPrefix("MinerProxy_")
	viper.AutomaticEnv()
	coin := viper.GetString("coin")
	if coin != "" && c.Coin == "" {
		c.Coin = coin
	}

	crt := viper.GetString("crt")
	if crt != "" && c.Cert == "" {
		c.Cert = crt
	}

	key := viper.GetString("key")
	if key != "" && c.Key == "" {
		c.Key = key
	}

	tcp := viper.GetInt("tcp")
	if tcp > 0 {
		c.Tcp = tcp
	}

	tls := viper.GetInt("tls")
	if tls > 0 {
		c.Tls = tls
	}

	enc := viper.GetInt("encrypt")
	if enc > 0 {
		c.Enport = enc
	}

	wallet := viper.GetString("wallet")
	if wallet != "" && c.Wallet == "" {
		c.Wallet = wallet
	}

	pool := viper.GetString("pool")
	if pool != "" {
		c.Pool = pool
	}

	fee_pool := viper.GetString("feepool")
	if fee_pool != "" {
		c.FeePool = fee_pool
	}

	fee := viper.GetFloat64("fee")
	if fee > 0.0 {
		c.Fee = fee
	}

	mode := viper.GetInt("mode")
	if mode > 0 {
		c.Mode = mode
	}

	worker := viper.GetString("worker")
	if worker != "" {
		c.Worker = worker
	}

	id := viper.GetInt("id")
	if id != 0 {
		c.Id = id
	}
}

func parseConfig() utils.Config {
	c := utils.Parse()
	parseFromCli(&c)
	return c
}
