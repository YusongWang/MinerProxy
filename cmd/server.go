package cmd

import (
	"miner_proxy/handles/eth"
	"miner_proxy/network"
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
		handle := eth.Handle{}
		s := serve.NewServe(net, &handle)
		s.StartLoop()
	},
}
