package main

import (
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"sync"
	"time"

	ipc "github.com/james-barrow/golang-ipc"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		//for {
		cc, err := ipc.StartClient(pool.ManageCmdPipeline, nil)
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		utils.Logger.Info("链接到manage成功")
		for {
			// for _, hand := range handle.Workers {
			// 	var json = jsoniter.ConfigCompatibleWithStandardLibrary
			// 	b, err := json.Marshal(hand)
			// 	if err != nil {
			// 		utils.Logger.Error(err.Error())
			// 		time.Sleep(time.Second * 10)
			// 		continue
			// 	}
			// 	err = cc.Write(1, b)
			// 	if err == nil {
			// 		utils.Logger.Info("发送成功!")
			// 	}
			// }
			err := cc.Write(111, []byte("Hello world!"))
			if err != nil {
				utils.Logger.Error(err.Error())
			}
			time.Sleep(time.Millisecond * 10)
		}
		//}
	}()

	wg.Wait()
}
