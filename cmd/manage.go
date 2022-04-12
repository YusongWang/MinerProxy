package cmd

import (
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"os"
	"os/exec"
	"sync"
	"time"

	ipc "github.com/james-barrow/golang-ipc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ManagerCmd)
}

var ManagerCmd = &cobra.Command{
	Use:   "manage",
	Short: "m",
	Long:  `manage`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		//TODO 解析SERVER配置文件。

		//TODO 解析Webconfig配置文件。

		//TODO 监听配置文件

		// 启动SERVER配置。

		// 启动web配置

		//Web Manage
		wg.Add(1)
		go Web(&wg)
		// TEST manage
		wg.Add(1)
		go Manage(&wg)

		wg.Wait()
	},
}

// 解析-配置文件。
func FristStart(wg *sync.WaitGroup) {
	sc, err := ipc.StartServer(pool.ManageCmdPipeline, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("Start Pipeline On " + pool.ManageCmdPipeline)

	for {
		msg, err := sc.Read()
		if err == nil {
			utils.Logger.Info("Server recieved: " + string(msg.Data))
		} else {
			utils.Logger.Error(err.Error())
			time.Sleep(time.Nanosecond * 10)
			break
		}
	}

	wg.Done()
}

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
			utils.Logger.Info("Server recieved: " + string(msg.Data))
		} else {
			utils.Logger.Error(err.Error())
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	wg.Done()
}

func Web(wg *sync.WaitGroup) {
web:
	web := exec.Command(os.Args[0], "web")
	err := web.Run()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	time.Sleep(time.Millisecond * 10)
	goto web
}
