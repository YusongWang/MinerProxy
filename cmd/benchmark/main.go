package main

import (
	"flag"
	"fmt"
	"miner_proxy/pack"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"os"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	pool   string
	thread int
	wallet string
)

func main() {

	flag.StringVar(&pool, "pool", "ssl://asia2.ethermine.org:5555", "set the pool sup tcp:// ssl://")
	flag.StringVar(&wallet, "wallet", "0xa324c686Cd081204F7A653E8435e18084AF81707", "Set the wallet address")
	flag.IntVar(&thread, "thread", 1000, "set thread num")
	if pool == "" {
		fmt.Println("Pool Is empty!")
		os.Exit(-1)
	}

	if wallet == "" {
		fmt.Println("Wallet Is empty!")
		os.Exit(-1)
	}

	if thread <= 0 {
		fmt.Println("Thread Is empty!")
		os.Exit(-1)
	}

	defer ants.Release()
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		worker := "P01"

		dev_job := &pack.Job{}
		dev_submit_job := make(chan []byte, 100)

		dev_pool, err := ethpool.New(pool, dev_job, dev_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}
		dev_pool.Login(wallet, worker)
		dev_pool.StartLoop()

		wg.Done()
	}

	for i := 0; i < thread; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}

	wg.Wait()
}
