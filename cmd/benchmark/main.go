package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"miner_proxy/global"
	"miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"os"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

var (
	pool   string
	thread int
	wallet string
)

func main() {

	flag.StringVar(&pool, "pool", "ssl://asia.etherminer.com:5555", "set the pool sup tcp:// ssl://")
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

		dev_job := &global.Job{}
		dev_submit_job := make(chan []byte, 100)

		dev_pool, err := ethpool.New(pool, dev_job, dev_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}
		dev_pool.Login(wallet, worker)
		wg.Add(1)
		go dev_pool.StartLoop()

		conn := *dev_pool.Conn
		for {
			if len(dev_job.Job) < 1 {
				continue
			}
			last_job := dev_job.Job[len(dev_job.Job)-1]
			submit := eth.ServerBaseReq{
				Id:     40,
				Method: "eth_submitWork",
				Params: last_job,
			}

			a, err := json.Marshal(submit)
			if err != nil {
				continue
			}
			a = append(a, '\n')
			conn.Write(a)
			time.Sleep(time.Minute * 1)
		}

		wg.Done()
	}

	for i := 0; i < thread; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}

	wg.Wait()
}
