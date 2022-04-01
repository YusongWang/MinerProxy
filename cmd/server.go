package cmd

import (
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	ethpack "miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/serve"
	"miner_proxy/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动MinerProxy核心，提供转发服务。",
	Long:  `无UI界面启动。`,
	Run: func(cmd *cobra.Command, args []string) {
		net, err := network.NewTcp(":38888")
		if err != nil {
			utils.Logger.Error("can't bind to addr", zap.String("端口", ":38888"))
		}

		dev_job := &ethpack.Job{}
		fee_job := &ethpack.Job{}

		dev_submit_job := make(chan []string, 10)
		fee_submit_job := make(chan []string, 10)
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

		handle := eth.Handle{}
		s := serve.NewServe(net, &handle)
		s.StartLoop()
	},
}
