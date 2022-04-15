package cmd

import (
	"fmt"
	"miner_proxy/global"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	routeRegister "miner_proxy/web/routes"
	"net/http"
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

		r := initRouter()
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

func initRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
	})

	routeRegister.RegisterApiRouter(router)

	// ReverseProxy
	// router.Use(proxy.ReverseProxy(map[string] string {
	// 	"localhost:4000" : "localhost:9090",
	// }))

	return router
}
