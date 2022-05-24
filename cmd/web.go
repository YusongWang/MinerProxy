package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"miner_proxy/global"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"miner_proxy/web/logics"
	"miner_proxy/web/models"
	routeRegister "miner_proxy/web/routes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ipc "github.com/james-barrow/golang-ipc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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

		go clacChart()

		// TODO： 将旷工5分钟内没有更新的旷工设置为离线。
		//go ChangeWorkerOffline()

		r := initRouter()
		utils.Logger.Info("Start Web Port On: " + strconv.Itoa(global.WebApp.Port))

		fmt.Println("Start Web Port On: " + strconv.Itoa(global.WebApp.Port) + "Password: " + global.WebApp.Password)

		r.Run(fmt.Sprintf(":%v", global.WebApp.Port))
	},
}

func FristStartIpcClients() {
	for _, app := range global.ManageApp.Config {
		// 逐一获得cmd执行任务。
		//fmt.Println("逐一获得cmd执行任务。")
		go StartIpcServer(app.ID)
		go ChangeWorkerOffline(app.ID)
	}
}

func ChangeWorkerOffline(id int) {
	for {
		for fullname, worker := range global.OnlinePools[id] {
			if !worker.IsOffline() {
				share_time_out := worker.ShareTime.Add(time.Minute * 5)
				now := time.Now()
				if now.Before(share_time_out) {
					worker.Logout()
					global.OnlinePools[id][fullname] = worker
				}
			}
		}

		time.Sleep(time.Second * 30)
	}
}

func StartIpcServer(id int) {
	pipename := pool.WebCmdPipeline + "_" + strconv.Itoa(id)
	log := utils.Logger.With(zap.String("IPC_NAME", pipename))
	config := ipc.ServerConfig{
		Encryption: true,
		MaxMsgSize: math.MaxInt,
	}

	for {
		sc, err := ipc.StartServer(pipename, &config)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Second * 60)
			continue
		}

		log.Info("IPC Server Ready to bind success")

		for {
			msg, err := sc.Read()
			if err != nil {
				log.Info("Server Error " + err.Error())
				//time.Sleep(time.Second * 30)
				break
			}

			if msg.MsgType <= 0 {
				continue
			}

			if msg.MsgType == 10 {
				err = sc.Write(10, []byte("PONG"))
				if err != nil {
					log.Error(err.Error())
				}
			}

			if msg.MsgType == 100 {
				var p map[string]global.Worker
				err := json.Unmarshal(msg.Data, &p)
				if err != nil {
					log.Error("格式化矿工状态失败", zap.String("data", string(msg.Data)))
					continue
				}
				global.OnlinePools[id] = p
			}
		}
	}
}

func StartIpcClient(id int) {
	pipename := pool.WebCmdPipeline + "_" + strconv.Itoa(id)
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
				var p map[string]global.Worker
				if msg.MsgType == 100 {
					err := json.Unmarshal(msg.Data, &p)
					if err != nil {
						log.Error("格式化矿工状态失败", zap.String("data", string(msg.Data)))
						continue
					}
					global.OnlinePools[id] = p

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
				time.Sleep(time.Second * 120)
			}
		}()

		wg.Wait()
	}
}

func initRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	routeRegister.RegisterApiRouter(router)
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
	//router.Handle("/", http.FileServer(AssetFile()))

	return router
}

//计算各个数据的图标
func clacChart() {
	for {
		//time.Sleep(time.Minute * 10)
		ethRes, etcRes := logics.ClacDashborad()

		insertTime := time.Now().Local().Unix()

		ethOnline := ethRes["online_worker"]
		ethOffline := ethRes["offline_worker"]
		ethHashrate := ethRes["total_hash"]

		var ethWorkerInfo models.WorkerChart
		ethWorkerInfo.Coin = "ETH"
		ethWorkerInfo.Time = insertTime
		if hash, ok := ethHashrate.(*big.Int); ok {
			ethWorkerInfo.Hashrate = hash
		}
		if offline, ok := ethOffline.(int); ok {
			ethWorkerInfo.Offline = offline
		}
		if online, ok := ethOnline.(int); ok {
			ethWorkerInfo.Online = online
		}
		err := models.InsertWorkerETH(ethWorkerInfo)
		if err != nil {
			utils.Logger.Info("insert Worker clac Chart Error" + err.Error())
		}

		etcOnline := etcRes["online_worker"]

		etcOffline := etcRes["offline_worker"]

		etcHashrate := etcRes["total_hash"]

		var etcWorkerInfo models.WorkerChart
		etcWorkerInfo.Coin = "ETC"
		etcWorkerInfo.Time = insertTime
		if hash, ok := etcHashrate.(*big.Int); ok {
			etcWorkerInfo.Hashrate = hash
		}

		if offline, ok := etcOffline.(int); ok {
			etcWorkerInfo.Offline = offline
		}

		if online, ok := etcOnline.(int); ok {
			etcWorkerInfo.Online = online
		}

		err = models.InsertWorkerETC(etcWorkerInfo)
		if err != nil {
			utils.Logger.Info("insert Worker clac Chart Error" + err.Error())
		}

		cpu := GetCpuPercent()
		mem := GetMemPercent()

		var sysinfo models.SystemChart
		sysinfo.Cpu = cpu
		sysinfo.Mem = mem
		sysinfo.Time = insertTime
		err = models.InsertSys(sysinfo)
		if err != nil {
			utils.Logger.Info("insert SystemInfo clac Chart Error" + err.Error())
		}

		time.Sleep(time.Minute * 10)
	}
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}
