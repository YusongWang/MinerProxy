package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动MinerProxy核心，提供转发服务。",
	Long:  `无UI界面启动。`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
