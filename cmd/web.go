package cmd

import (
	"fmt"
	"miner_proxy/global"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	ipc "github.com/james-barrow/golang-ipc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	WebCmd.Flags().String("password", "admin123", "指定web密码")
	viper.BindPFlag("password", WebCmd.Flags().Lookup("password"))

	WebCmd.Flags().Int("port", 9898, "指定web端口")
	viper.BindPFlag("port", WebCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(WebCmd)
}

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "w",
	Long:  `web`,
	Run: func(cmd *cobra.Command, args []string) {
		go StartIpcServer()

		port := viper.GetInt("port")
		global.WebApp.Port = port
		password := viper.GetString("password")
		global.WebApp.Password = password

		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		utils.Logger.Info("Start Web Port On: " + strconv.Itoa(global.WebApp.Port))

		r.Run(fmt.Sprintf(":%v", global.WebApp.Port))
	},
}

func StartIpcServer() {
	sc, err := ipc.StartServer(pool.WebCmdPipeline, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("Start Web Pipeline On: " + pool.WebCmdPipeline)

	for {
		msg, err := sc.Read()
		if err == nil {
			utils.Logger.Info("Server recieved: "+string(msg.Data), zap.Int("type", msg.MsgType))
		} else {
			utils.Logger.Error(err.Error())
			break
		}
	}
}
