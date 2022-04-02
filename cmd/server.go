package cmd

import (
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	ethpack "miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/serve"
	"miner_proxy/utils"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("coin", "ETH", "指定需要代理的币种")
	viper.BindPFlag("coin", serverCmd.Flags().Lookup("coin"))

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

	serverCmd.Flags().String("pool", "", "指定需要配置代理中转的矿池地址\n格式: (ssl://)tcp://asia2.ethermine.org:4444")
	viper.BindPFlag("pool", serverCmd.Flags().Lookup("pool"))

	serverCmd.Flags().String("feepool", "", "指定抽水的矿池地址\n格式： (ssl://)tcp://asia2.ethermine.org:4444")
	viper.BindPFlag("fee_pool", serverCmd.Flags().Lookup("fee_pool"))

	serverCmd.Flags().Int("mode", 0, "中转模式: 1. 直连 2. 抽水")
	viper.BindPFlag("mode", serverCmd.Flags().Lookup("mode"))

	serverCmd.Flags().Float64("fee", 0.0, "抽水率(%)默认 2% 支持一位小数点。")
	viper.BindPFlag("fee", serverCmd.Flags().Lookup("fee"))
	// serverCmd.Flags().String("config", "./config.yaml", "指定配置文件")
	// viper.BindPFlag("config", serverCmd.Flags().Lookup("config"))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动MinerProxy核心，提供转发服务。",
	Long:  `无UI界面启动。`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(pflag.CommandLine)
		config := parseConfig()

		if err := config.Check(); err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}

		if config.Mode == 1 {
			switch config.Coin {
			case utils.ETH:

			default:
				utils.Logger.Error("暂未支持的币种")
				os.Exit(99)
			}
		} else if config.Mode == 2 {
			switch config.Coin {
			case utils.ETH:

			default:
				utils.Logger.Error("暂未支持的币种")
				os.Exit(99)
			}
		} else {
			utils.Logger.Error("不支持的Mode参数")
			os.Exit(99)
		}

		dev_job := &ethpack.Job{}
		fee_job := &ethpack.Job{}

		dev_submit_job := make(chan []string, 10)
		fee_submit_job := make(chan []string, 10)

		net, err := network.NewTcp(":38880")
		if err != nil {
			utils.Logger.Error("can't bind to TCP addr", zap.String("端口", ":38880"))
			os.Exit(99)
			return
		}

		nettls, err := network.NewTls("cert.pem", "key.pem", ":38888")
		if err != nil {
			utils.Logger.Error("can't bind to SSL addr", zap.String("端口", ":38888"))
			os.Exit(99)
			return
		}

		// 开启两个抽水线程

		// 开发者线程
		dev_pool, err := ethpool.NewEthStratumServerSsl("api.wangyusong.com:8443", dev_job, dev_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		dev_pool.Login("0x3602b50d3086edefcd9318bcceb6389004fb14ee")
		go dev_pool.StartLoop()

		// 中转线程
		fee_pool, err := ethpool.NewEthStratumServerSsl("api.wangyusong.com:8443", fee_job, fee_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		fee_pool.Login("0x3602b50d3086edefcd9318bcceb6389004fb14ee")
		go fee_pool.StartLoop()

		var wg sync.WaitGroup
		handle := eth.Handle{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			s := serve.NewServe(net, &handle)
			s.StartLoop()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			s := serve.NewServe(nettls, &handle)
			s.StartLoop()
		}()

		wg.Wait()
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
	if pool != "" && c.Pool == "" {
		c.Pool = pool
	}

	fee_pool := viper.GetString("feepool")
	if fee_pool != "" && c.FeePool == "" {
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

}

func parseConfig() utils.Config {
	c := utils.Parse()
	parseFromCli(&c)
	return c
}
