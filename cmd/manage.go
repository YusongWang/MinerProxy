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
		//Web Manage
		wg.Add(1)
		go Web(&wg)
		// TEST manage
		wg.Add(1)
		go Manage(&wg)

		wg.Wait()
	},
}

func Manage(wg *sync.WaitGroup) {
	sc, err := ipc.StartServer(pool.ManageCmdPipeline, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	for {
		msg, err := sc.Read()
		if err == nil {
			utils.Logger.Info("Server recieved: " + string(msg.Data))
		} else {
			utils.Logger.Error(err.Error())
			time.Sleep(time.Second)
			break
		}
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
	goto web

	wg.Done()
}
