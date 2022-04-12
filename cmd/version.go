package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印当前版本号",
	Long:  `打印当前版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本号: %v 当前commit: %v \n", version, commit)
	},
}
