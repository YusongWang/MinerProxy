package cmd

import (
	"fmt"
	"miner_proxy/global"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印当前版本号",
	Long:  `打印当前版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本号: %v 当前commit: %v \n", global.Version, global.Commit)
	},
}
