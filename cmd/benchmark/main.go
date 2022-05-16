package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
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

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
func main() {

	flag.StringVar(&pool, "pool", "tcp://localhost:8812", "set the pool sup tcp:// ssl://")
	flag.StringVar(&wallet, "wallet", "0xa324c686Cd081204F7A653E8435e18084AF81707", "Set the wallet address")
	flag.IntVar(&thread, "thread", 10, "set thread num")
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
	// 增大文件描述符上限
	utils.IncreaseFDLimit()

	defer ants.Release()
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		worker := CreateRandomString(10)

		var dev_job []global.Job
		dev_submit_job := make(chan []byte, 100)

		dev_pool, err := ethpool.New(pool, &dev_job, dev_submit_job)
		if err != nil {
			utils.Logger.Error(err.Error())
			os.Exit(99)
		}
		dev_pool.Login(wallet, worker)
		wg.Add(1)
		go dev_pool.StartLoop()

		conn := *dev_pool.Conn
		for {
			if len(dev_job) < 1 {
				continue
			}

			last_job := dev_job[len(dev_job)-1]
			submit := eth.ServerBaseReq{
				Id:     40,
				Method: "eth_submitWork",
				Params: []string{last_job.Target, last_job.JobId, last_job.Diff},
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
