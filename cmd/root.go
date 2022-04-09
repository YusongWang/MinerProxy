package cmd

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "MinerProxy",
	Short: "高性能矿工代理工具",
	Long:  `提供高性能的矿工服务转发`,
	Run: func(cmd *cobra.Command, args []string) {
		//gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		r.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
}
