package cmd

import (
	"encoding/json"
	"fmt"
	"miner_proxy/global"
	"miner_proxy/pack"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	routeRegister "miner_proxy/web/routes"
	"net/http"
	"strconv"
	"sync"
	"time"

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
		web_notify_ch := make(chan int)
		proxy_notify_ch := make(chan int)

		// 解析SERVER配置文件。
		// 监听配置文件
		InitializeConfig(web_notify_ch, proxy_notify_ch)
		FristStartIpcClients()

		port := viper.GetInt("port")
		global.WebApp.Port = port
		password := viper.GetString("password")
		global.WebApp.Password = password

		r := initRouter()
		utils.Logger.Info("Start Web Port On: " + strconv.Itoa(global.WebApp.Port))
		r.Run(fmt.Sprintf(":%v", global.WebApp.Port))
	},
}

func FristStartIpcClients() {
	for _, app := range global.ManageApp.Config {
		// 逐一获得cmd执行任务。
		fmt.Println("逐一获得cmd执行任务。")
		go StartIpcClient(app.ID)
	}
}

func StartIpcClient(id int) {
	pipename := pool.WebCmdPipeline + "_" + strconv.Itoa(id)
	//pipename := pool.WebCmdPipeline
	log := utils.Logger.With(zap.String("IPC_NAME", pipename))
	for {
		cc, err := ipc.StartClient(pipename, nil)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Second * 60)
			continue
		}
		log.Info("IPC client Ready to Connect!")

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				msg, err := cc.Read()
				if err != nil {
					log.Info("Ipc Channel Close")
				}
				var p map[string]pack.Worker
				if msg.MsgType == 100 {
					err := json.Unmarshal(msg.Data, &p)
					if err != nil {
						log.Error("格式化矿工状态失败", zap.String("data", string(msg.Data)))
						continue
					}
					global.OnlinePools[id] = p
					//log.Info("Web 收到矿工信息", zap.Any("pool_workers", OnlinePools))
					continue
				}
				log.Info("Web recieved: "+string(msg.Data), zap.Int("type", msg.MsgType))
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err = cc.Write(111, []byte("hello"))
				if err != nil {
					log.Error(err.Error())
				}
				time.Sleep(time.Second * 10)
			}
		}()

		wg.Wait()
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

	return router
}
