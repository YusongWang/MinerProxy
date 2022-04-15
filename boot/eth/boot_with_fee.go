package eth

import (
	"fmt"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"miner_proxy/pack"
	"time"

	pool "miner_proxy/pools"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/serve"
	"miner_proxy/utils"

	"os"
	"sync"

	ipc "github.com/james-barrow/golang-ipc"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func BootWithFee(c utils.Config) error {

	dev_job := &pack.Job{}
	fee_job := &pack.Job{}

	dev_submit_job := make(chan []byte, 100)
	fee_submit_job := make(chan []byte, 100)

	// 中转线程
	fee_pool, err := ethpool.New(c.Feepool, fee_job, fee_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}

	fee_pool.Login(c.Wallet, c.Worker)
	go fee_pool.StartLoop()

	// 开发者线程
	dev_pool, err := ethpool.New(pool.ETH_POOL, dev_job, dev_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}

	dev_pool.Login(pool.ETH_WALLET, "devfee0.0.1")
	go dev_pool.StartLoop()

	var wg sync.WaitGroup
	handle := eth.Handle{
		Devjob:  dev_job,
		Feejob:  fee_job,
		DevConn: &dev_pool.Conn,
		FeeConn: &fee_pool.Conn,
		SubDev:  &dev_submit_job,
		SubFee:  &fee_submit_job,
		Workers: make(map[string]*pack.Worker),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			cc, err := ipc.StartClient(pool.WebCmdPipeline, nil)
			if err != nil {
				utils.Logger.Error(err.Error())
				time.Sleep(time.Second * 60)
				continue
			}

			go func() {
				for {
					msg, err := cc.Read()
					if err != nil {
						utils.Logger.Info("Ipc Channel Close")
					}
					utils.Logger.Info("Web -> Proxy ", zap.Any("msg", msg))
					time.Sleep(time.Second * 10)
				}
			}()

			utils.Logger.Info("链接到manage成功")
			for {
				for _, hand := range handle.Workers {
					var json = jsoniter.ConfigCompatibleWithStandardLibrary
					b, err := json.Marshal(hand)
					if err != nil {
						utils.Logger.Error(err.Error())
						//time.Sleep(time.Second * 60)
						continue
					}
					utils.Logger.Info("写入Worker信息!", zap.Any("worker", hand))
					err = cc.Write(100, b)
					if err == nil {
						utils.Logger.Info("发送成功!", zap.Any("worker", hand))
					} else {
						utils.Logger.Info("发送失败!")
						utils.Logger.Error(err.Error())
					}
				}
				time.Sleep(time.Second * 60)
			}
		}
	}()

	utils.Logger.Info("Start the Server And ready To serve")

	if c.TCP > 0 {
		port := fmt.Sprintf(":%v", c.TCP)
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

	if c.TLS > 0 {
		port := fmt.Sprintf(":%v", c.TLS)
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
