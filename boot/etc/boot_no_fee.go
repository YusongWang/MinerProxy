package etc

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
	handle := eth.NoFeeHandle{}

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
