package eth

import (
	"fmt"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"miner_proxy/serve"
	"miner_proxy/utils"

	"os"
	"sync"

	"go.uber.org/zap"
)

func BootNoFee(c utils.Config) error {
	// wait
	var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	StartIpcClient(c.ID)
	// }()

	handle := eth.NoFeeHandle{}

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
