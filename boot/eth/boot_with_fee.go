package eth

import (
	"fmt"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"miner_proxy/pack"
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

func StartIpcServer(id int, handle *eth.Handle) {
	pipename := pool.WebCmdPipeline + ":" + strconv.Itoa(id)
	log := utils.Logger.With(zap.String("IPC_NAME", pipename))

	for {
		sc, err := ipc.StartServer(pool.WebCmdPipeline, nil)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("Start Proxy To Web Pipeline On: " + pipename)

		go func() {
			for {
				msg, err := sc.Read()
				if err == nil {
					log.Info("Server recieved: "+string(msg.Data), zap.Int("type", msg.MsgType))
					return
				} else {
					log.Error(err.Error())
					break
				}
			}
		}()

		for {
			for _, hand := range handle.Workers {
				var json = jsoniter.ConfigCompatibleWithStandardLibrary
				b, err := json.Marshal(hand)
				if err != nil {
					log.Error(err.Error())
					//time.Sleep(time.Second * 60)
					continue
				}
				log.Info("写入Worker信息!", zap.Any("worker", hand))
				err = sc.Write(100, b)
				if err == nil {
					log.Info("发送成功!", zap.Any("worker", hand))
				} else {
					log.Info("发送失败!")
					log.Error(err.Error())
				}
			}
			time.Sleep(time.Second * 60)
		}
	}
}

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
		StartIpcServer(c.ID, &handle)
		// for {

		// 	pipename := pool.WebCmdPipeline + ":" + strconv.Itoa(c.ID)
		// 	log := utils.Logger.With(zap.String("IPC_NAME", pipename))
		// 	cc, err := ipc.StartClient(pipename, nil)
		// 	if err != nil {
		// 		log.Error(err.Error())
		// 		time.Sleep(time.Second * 60)
		// 		continue
		// 	}

		// 	go func() {
		// 		for {
		// 			msg, err := cc.Read()
		// 			if err != nil {
		// 				log.Info("Ipc Channel Close")
		// 			}
		// 			if msg.MsgType == -1 {
		// 				//TODO
		// 				continue
		// 			}
		// 			log.Info("Proxy ->  Web", zap.Any("msg", msg))
		// 			time.Sleep(time.Second * 10)
		// 		}
		// 	}()

		// 	for {
		// 		for _, hand := range handle.Workers {
		// 			var json = jsoniter.ConfigCompatibleWithStandardLibrary
		// 			b, err := json.Marshal(hand)
		// 			if err != nil {
		// 				log.Error(err.Error())
		// 				//time.Sleep(time.Second * 60)
		// 				continue
		// 			}
		// 			log.Info("写入Worker信息!", zap.Any("worker", hand))
		// 			err = cc.Write(100, b)
		// 			if err == nil {
		// 				log.Info("发送成功!", zap.Any("worker", hand))
		// 			} else {
		// 				log.Info("发送失败!")
		// 				log.Error(err.Error())
		// 			}
		// 		}
		// 		time.Sleep(time.Second * 60)
		// 	}
		// }
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
