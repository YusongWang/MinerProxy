package cmd

import (
	"fmt"
	"miner_proxy/utils"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// 解析SERVER配置文件。
		// 监听配置文件
		InitializeConfig(web_notify_ch, proxy_notify_ch)

		// 启动web配置
		wg.Add(1)
		go Web(&wg, web_notify_ch)

		// 启动代理watchdog
		wg.Add(1)
		go Proxy(&wg, proxy_notify_ch)

		// 等待退出（永远不会退出）
		wg.Wait()
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

// func Manage(wg *sync.WaitGroup) {
// 	sc, err := ipc.StartServer(pool.ManageCmdPipeline, nil)
// 	if err != nil {
// 		utils.Logger.Error(err.Error())
// 		return
// 	}

// 	utils.Logger.Info("Start Pipeline On " + pool.ManageCmdPipeline)

// 	for {
// 		msg, err := sc.Read()
// 		if err == nil {
// 			utils.Logger.Info("Server recieved: "+string(msg.Data), zap.Int("type", msg.MsgType))
// 		} else {
// 			utils.Logger.Error(err.Error())
// 			break
// 		}
// 	}

// 	wg.Done()
// }

func Web(wg *sync.WaitGroup, restart chan int) {
web:
	fmt.Println(os.Args[0], "web", "--port", strconv.Itoa(ManageApp.Web.Port), "--password", ManageApp.Web.Password)
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

	time.Sleep(time.Second * 10)
	goto web
}

func Proxy(wg *sync.WaitGroup, restart chan int) {
	FristStart()
proxy:
	//TODO 启动所有proxy_worker
	// 注册为一个临时数组、管理所有worker. id 为当前结构注册的 ID
	//func() {
	for {
		select {
		case id := <-restart:
			utils.Logger.Info("重启代理ID: " + strconv.Itoa(id))
			//FIXME 处理旧任务？ 如果任务ID 变更旧任务就要删掉。

			// for online_id, cmd := range ManagePool.Online {
			// 	for _, app := range ManageApp.Config {
			// 		if app.ID == online_id {
			// 			ProcessProxy(&app)
			// 		}
			// 	}
			// }

			if ManagePool.Online[id] == nil {
				for _, app := range ManageApp.Config {
					if app.ID == id {
						ProcessProxy(app)
					}
				}
			} else {
				ManagePool.Online[id].Process.Kill()
			}
		}
	}
	//}()

	// 注册一个chan 接收ID作为重启。如果这个ID不在数组中就新增一个代理池
	time.Sleep(time.Second * 10)
	goto proxy
}

func FristStart() {
	for _, app := range ManageApp.Config {
		// 逐一获得cmd执行任务。
		fmt.Println("逐一获得cmd执行任务。")
		fmt.Println(app)
		go ProcessProxy(app)
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
	v.AddConfigPath("/etc/appname/")  // path to look for the config file in
	v.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	v.AddConfigPath(".")              // optionally look for config in the working directory
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		utils.Logger.Error(err.Error())
		return v
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		utils.Logger.Info("config file changed:" + in.Name)

		//copy(ManageApp, conf)
		conf := *ManageApp
		//conf := *ManageApp

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
			is_new := true
			//FIXME 如果这里为空可能不会新增代理
			for _, old_app := range conf.Config {
				if app.ID == old_app.ID {
					if checkConfigChange(old_app, app) {
						is_new = false
						proxy_restart <- app.ID
					}
				}
			}
			if is_new {
				proxy_restart <- app.ID
			}
		}
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

func ProcessProxy(c utils.Config) {
proxy:
	//--coin ETH --tcp 38888 --pool tcp://asia2.ethermine.org:4444 --feepool
	//tcp://asia2.ethermine.org:4444
	//--mode 2 --wallet 0x3602b50d3086edefcd9318bcceb6389004fb14ee --fee 5
	fmt.Println(os.Args[0],
		"server",
		"--id",
		strconv.Itoa(c.ID),
		"--coin",
		c.Coin,
		"--tcp",
		strconv.Itoa(c.TCP),
		"--tls",
		strconv.Itoa(c.TLS),
		"--encrypt",
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
		"--crt",
		c.Cert,
		"--online")
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
		"--encrypt",
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
		"--crt",
		c.Cert,
		"--online",
	)
	if !c.Online {
		return
	}

	fmt.Println(c)

	ManagePool.Online[c.ID] = p
	utils.Logger.Info("启动代理软件")

	err := p.Run()
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	time.Sleep(time.Second * 10)
	goto proxy
}
