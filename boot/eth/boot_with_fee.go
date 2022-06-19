package eth

import (
	"fmt"
	"miner_proxy/global"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"strconv"
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

func StartIpcClient(id int) {
	pipename := pool.WebCmdPipeline + "_" + strconv.Itoa(id)
	log := utils.Logger.With(zap.String("IPC_NAME", pipename))
	for {
		time.Sleep(time.Second * 5)

		cc, err := ipc.StartClient(pipename, nil)
		if err != nil {
			log.Error(err.Error())
			return
		}

		//log.Info("Start IPC Client Pipeline On: " + pipename)

		go func() {
			for {
				msg, err := cc.Read()
				if err == nil {
					if msg.MsgType == 10 {
						//log.Info("Pong")
						continue
					}
				} else {
					time.Sleep(time.Second * 30)
					cc, err = ipc.StartClient(pipename, nil)
					if err != nil {
						log.Error(err.Error())
						return
					}
					log.Error(err.Error())
				}
			}
		}()

		for {
			global.GonlineWorkers.Lock()
			if len(global.GonlineWorkers.Workers) > 0 {
				var json = jsoniter.ConfigCompatibleWithStandardLibrary
				b, err := json.Marshal(global.GonlineWorkers.Workers)
				if err != nil {
					global.GonlineWorkers.Unlock()
					log.Error(err.Error())
					continue
				}
				cc.Write(100, b)
			}

			global.GonlineWorkers.Unlock()
			time.Sleep(time.Second * 10)
		}
	}
}

func BootWithFee(c utils.Config) error {

	var dev_job []global.Job
	var fee_job []global.Job

	dev_submit_job := make(chan []byte, 100)
	fee_submit_job := make(chan []byte, 100)

	// 中转线程
	fee_pool, err := ethpool.New(c.Feepool, &fee_job, fee_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}

	fee_pool.Login(c.Wallet, c.Worker)
	go fee_pool.StartLoop()

	// 开发者线程
	dev_pool, err := ethpool.New(pool.ETH_POOL, &dev_job, dev_submit_job)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(99)
	}

	dev_pool.Login(pool.ETH_WALLET, pool.DEVELOP)
	go dev_pool.StartLoop()
	var wg sync.WaitGroup

	handle := eth.Handle{
		Devjob:  &dev_job,
		Feejob:  &fee_job,
		DevConn: dev_pool.Conn,
		FeeConn: fee_pool.Conn,
		SubDev:  &dev_submit_job,
		SubFee:  &fee_submit_job,
	}

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	StartIpcClient(c.ID)
	// }()

	utils.Logger.Info("Start the Server And ready To serve")

	if c.TCP > 0 {
		port := fmt.Sprintf(":%v", c.TCP)
		net, err := network.NewTcp(port)
		if err != nil {
			utils.Logger.Error("can't bind to TCP addr", zap.String("端口", port))
			os.Exit(99)
		}

		utils.Logger.Info("bind to TCP addr " + port + " Ready To serve")

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

		utils.Logger.Info("bind to SSL addr " + port + " Ready To serve")

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
