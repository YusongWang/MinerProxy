package cmd

import (
	"miner_proxy/utils"
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
	Short: "启动管理端",
	Long:  `启动管理端`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(1)
		go Manage(&wg)

		wg.Wait()
	},
}

func Manage(wg *sync.WaitGroup) {
	sc, err := ipc.StartServer("MinerProxy", nil)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	for {
		msg, err := sc.Read()
		if err == nil {
			utils.Logger.Info("Server recieved: " + string(msg.Data) + " - Message type: ")
		} else {
			utils.Logger.Error(err.Error())
			time.Sleep(time.Second * 10)
			break
		}
	}

	wg.Done()
}
