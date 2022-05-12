package test

import (
	"bufio"
	"io"
	"miner_proxy/global"
	"miner_proxy/pack"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Test struct {
	log *zap.Logger
}

func (hand *Test) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
	addr string,
	worker *pack.Worker,
) (io.ReadWriteCloser, error) {
	hand.log.Info("On Miner Connect To Pool " + config.Pool)
	pool, err := utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}

	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			hand.log.Info("矿池: " + string(buf))
			c.Write(buf)
		}
	}()

	return pool, nil
}

func (hand *Test) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
	data *[]byte,
	worker *pack.Worker,
) (out []byte, err error) {
	hand.log.Info("矿机: " + string(*data))
	(*pool).Write(*data)

	out = nil
	err = nil
	return
}

func (hand *Test) OnClose(worker *pack.Worker) {
	hand.log.Info("OnClose !!!!!")
}

func (hand *Test) SetLog(log *zap.Logger) {
	hand.log = log
}
