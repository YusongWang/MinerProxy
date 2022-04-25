package main

import (
	"miner_proxy/pack"
	pool "miner_proxy/pools"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"os"
	"sync"

	"github.com/panjf2000/ants/v2"
)

func main() {
	defer ants.Release()
	runTimes := 1000
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		worker := "Hello world"

		dev_job := &pack.Job{}
		dev_submit_job := make(chan []byte, 100)

		dev_pool, err := ethpool.New("ssl://api.wangyusong.com:8443", dev_job, dev_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}
		dev_pool.Login(pool.ETH_WALLET, worker)
		dev_pool.StartLoop()

		wg.Done()
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}

	wg.Wait()
}
