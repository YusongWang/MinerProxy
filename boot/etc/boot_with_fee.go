package etc

import (
	"fmt"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"miner_proxy/pack"
	pool "miner_proxy/pools"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/serve"
	"miner_proxy/utils"

	"os"
	"sync"

	"go.uber.org/zap"
)

func BootWithFee(c utils.Config) error {
	dev_job := &pack.Job{}
	fee_job := &pack.Job{}

	dev_submit_job := make(chan []byte, 100)
	fee_submit_job := make(chan []byte, 100)
	// 中转线程
	dev_pool, err := ethpool.New(c.FeePool, fee_job, fee_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}

	//TODO check wallet len and Start with 0x
	dev_pool.Login(c.Wallet, c.Worker)
	go dev_pool.StartLoop()

	// 开发者线程
	fee_pool, err := ethpool.New(pool.ETC_POOL, dev_job, dev_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}
	fee_pool.Login(pool.ETC_WALLET, "DEVFEE-0.1")
	go fee_pool.StartLoop()

	// wait
	var wg sync.WaitGroup
	handle := eth.Handle{
		Devjob:  dev_job,
		Feejob:  fee_job,
		DevConn: &dev_pool.Conn,
		FeeConn: &fee_pool.Conn,
	}

	utils.Logger.Info("Start the Server And ready To serve")

	if c.Tcp > 0 {
		port := fmt.Sprintf(":%v", c.Tcp)
		net, err := network.NewTcp(port)
		if err != nil {
			utils.Logger.Error("can't bind to TCP addr", zap.String("端口", port))
			os.Exit(99)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			s := serve.NewServe(net, &handle, &c)
			s.StartLoop()
		}()
	}

	if c.Tls > 0 {
		port := fmt.Sprintf(":%v", c.Tls)
		nettls, err := network.NewTls(c.Cert, c.Key, port)
		if err != nil {
			utils.Logger.Error("can't bind to SSL addr", zap.String("端口", port))
			os.Exit(99)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			s := serve.NewServe(nettls, &handle, &c)
			s.StartLoop()
		}()
	}

	wg.Wait()
	return nil
}
