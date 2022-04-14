package cmd

import (
	"fmt"
	"log"
	"miner_proxy/global"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	ipc "github.com/james-barrow/golang-ipc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ManageConfig struct {
	Config []utils.Config `json:"config"`
	Web    struct {
		Port     int    `json:"port"`
		Password string `json:"password"`
	} `json:"web"`
}

var ManageApp = new(ManageConfig)

type PoolConfig struct {
	Online []*exec.Cmd
}

var ManagePool PoolConfig

// func Init() {

// }

var rootCmd = &cobra.Command{
	Use:   "MinerProxy",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		web_notify_ch := make(chan int)
		proxy_notify_ch := make(chan int)
		// deamon the watch dog.
		//for i := 0; i < 1000; i++ {
		ManagePool.Online = make([]*exec.Cmd, 1000)
		//}
		// viper watch the File change. Web save the pool list

		// web select the Pool and customer setting pool

		// start Parse the web strings

		// if not set on web set the password and web port . gen the config,and restart self

		//gin.SetMode(gin.ReleaseMode)
		var wg sync.WaitGroup
		//TODO 解析SERVER配置文件。

		//TODO 解析Webconfig配置文件。

		//TODO 监听配置文件
		InitializeConfig(web_notify_ch, proxy_notify_ch)
		fmt.Println(global.WebApp)

		//ManagePool.Online[0] =
		// 启动SERVER配置。

		// 启动web配置

		//Web Manage
		wg.Add(1)
		go Web(&wg, web_notify_ch)
		wg.Add(1)
		go Proxy(&wg, proxy_notify_ch)
		// TEST manage
		wg.Add(1)
		go Manage(&wg)
		cc, err := ipc.StartClient(pool.ManageCmdPipeline, nil)
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		go func() {
			for {
				m, err := cc.Read()

				if err != nil {
					// An error is only returned if the recieved channel has been closed,
					//so you know the connection has either been intentionally closed or has timmed out waiting to connect/re-connect.
					break
				}

				if m.MsgType == -1 { // message type -1 is status change
					log.Println("Status: " + m.Status)
				}

				if m.MsgType == -2 { // message type -2 is an error, these won't automatically cause the recieve channel to close.
					log.Println("Error: " + err.Error())
				}

				if m.MsgType > 0 { // all message types above 0 have been recieved over the connection

					log.Println("Message type: ", m.MsgType)
					log.Println("Client recieved: " + string(m.Data))
				}
			}
		}()

		for {
			cc.Write(20, []byte("hello world"))
			time.Sleep(time.Second * 30)
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
}

// 解析-配置文件。
// func FristStart(wg *sync.WaitGroup) {
// 	sc, err := ipc.StartServer(pool.ManageCmdPipeline, nil)
// 	if err != nil {
// 		utils.Logger.Error(err.Error())
// 		return
// 	}

// 	utils.Logger.Info("Start Pipeline On " + pool.ManageCmdPipeline)

// 	for {
// 		msg, err := sc.Read()
// 		if err == nil {
// 			utils.Logger.Info("Server recieved: " + string(msg.Data))
// 		} else {
// 			utils.Logger.Error(err.Error())
// 			time.Sleep(time.Nanosecond * 10)
// 			break
// 		}
// 	}
// 	utils.Logger.Info("IPC exit()")
// 	wg.Done()
// }

func Manage(wg *sync.WaitGroup) {
	sc, err := ipc.StartServer(pool.ManageCmdPipeline, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("Start Pipeline On " + pool.ManageCmdPipeline)

	for {
		msg, err := sc.Read()
		if err == nil {
			utils.Logger.Info("Server recieved: "+string(msg.Data), zap.Int("type", msg.MsgType))
		} else {
			utils.Logger.Error(err.Error())
			break
		}
	}

	wg.Done()
}

func Web(wg *sync.WaitGroup, restart chan int) {
web:
	web := exec.Command(os.Args[0], "web", "--port", strconv.Itoa(ManageApp.Web.Port), "--password", ManageApp.Web.Password)
	go func() {
		<-restart
		fmt.Println("收到重启命令 Kill")
		web.Process.Kill()
	}()

	err := web.Run()
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	time.Sleep(time.Millisecond * 10)
	goto web
}

func Proxy(wg *sync.WaitGroup, restart chan int) {
	FristStart()
proxy:
	//TODO 启动所有proxy_worker
	// 注册为一个临时数组、管理所有worker. id 为当前结构注册的 ID
	go func() {
		for {
			select {
			case id := <-restart:
				utils.Logger.Info("重启代理ID: " + strconv.Itoa(id))
				ManagePool.Online[id].Process.Kill()
			}
		}
	}()
	// 注册一个chan 接收ID作为重启。如果这个ID不在数组中就新增一个代理池

	time.Sleep(time.Millisecond * 10)
	goto proxy
}

func FristStart() {
	for _, app := range ManageApp.Config {
		// 逐一获得cmd执行任务。
		ProcessProxy(&app)
	}
}

func InitializeConfig(web_restart chan int, proxy_restart chan int) *viper.Viper {
	// 设置配置文件路径
	config := "config.json"
	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("MINER_CONFIG"); configEnv != "" {
		config = configEnv
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		utils.Logger.Error(err.Error())
		return v
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		utils.Logger.Info("config file changed:" + in.Name)

		// 保存旧配置。
		conf := *ManageApp

		// Web 重载配置
		if err := v.Unmarshal(&ManageApp); err != nil {
			utils.Logger.Error(err.Error())
		}

		if ManageApp.Web.Password != conf.Web.Password || ManageApp.Web.Port != conf.Web.Port {
			//notify web
			web_restart <- 1
			fmt.Println("发送重启命令1")
		}

		// 检查 proxy 是否重启。
		for _, app := range ManageApp.Config {
			//FIXME 如果这里为空可能不会新增代理
			c := conf.Config[app.ID]
			if checkConfigChange(c, app) {
				proxy_restart <- app.ID
			}
		}

		fmt.Println(ManageApp)
	})

	// 将配置赋值给全局变量
	if err := v.Unmarshal(&ManageApp); err != nil {
		utils.Logger.Error(err.Error())
	}

	fmt.Println(ManageApp)
	return v
}

func checkConfigChange(old, new utils.Config) bool {
	if old.Cert != new.Cert {
		return true
	}
	if old.ID != new.ID {
		return true
	}
	if old.TCP != new.TCP {
		return true
	}
	if old.Coin != new.Coin {
		return true
	}
	if old.TLS != new.TLS {
		return true
	}
	if old.Enport != new.Enport {
		return true
	}
	if old.Key != new.Key {
		return true
	}
	if old.Pool != new.Pool {
		return true
	}

	if old.Feepool != new.Feepool {
		return true
	}

	if old.Fee != new.Fee {
		return true
	}

	if old.Worker != new.Worker {
		return true
	}

	if old.Wallet != new.Wallet {
		return true
	}

	if old.Mode != new.Mode {
		return true
	}

	if old.Online != new.Online {
		return true
	}

	return false
}

func ProcessProxy(c *utils.Config) {
proxy:
	//--coin ETH --tcp 38888 --pool tcp://asia2.ethermine.org:4444 --feepool
	//tcp://asia2.ethermine.org:4444
	//--mode 2 --wallet 0x3602b50d3086edefcd9318bcceb6389004fb14ee --fee 5
	p := exec.Command(
		os.Args[0],
		"server",
		"--id",
		strconv.Itoa(c.ID),
		"--coin",
		c.Coin,
		"--tcp",
		strconv.Itoa(c.TCP),
		"--tls",
		strconv.Itoa(c.TLS),
		"--enport",
		strconv.Itoa(c.Enport),
		"--pool",
		c.Pool,
		"--feepool",
		c.Feepool,
		"--fee",
		fmt.Sprintf("%f", c.Fee),
		"--mode",
		strconv.Itoa(c.Mode),
		"--wallet",
		c.Wallet,
		"--worker",
		c.Worker,
		"--key",
		c.Key,
		"--cert",
		c.Cert,
	)
	ManagePool.Online[c.ID] = p
	err := p.Run()
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	time.Sleep(time.Second * 10)
	goto proxy
}
