package eth_stratum

import (
	"io"
	"miner_proxy/fee"
	"miner_proxy/pack"
	ethpack "miner_proxy/pack/eth_stratum"
	"miner_proxy/utils"
	"strings"

	"bufio"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Handle struct {
	log     *zap.Logger
	Devjob  *pack.Job
	Feejob  *pack.Job
	DevConn *io.ReadWriteCloser
	FeeConn *io.ReadWriteCloser
}

func (hand *Handle) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *fee.Fee,
	addr string,
	id *string,
) (io.ReadWriteCloser, error) {
	hand.log.Info("On Miner Connect To Pool " + config.Pool)
	pool, err := utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		//writer := bufio.NewWriter(c)
		//log := hand.log.With(zap.String("Miner", c.RemoteAddr().String()))
		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}

			var push ethpack.JSONPushMessage
			if err = json.Unmarshal([]byte(buf), &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result {
						// TODO
						hand.log.Info("有效份额", zap.Any("RPC", string(buf)))
					} else {
						hand.log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else if _, ok := push.Result.([]interface{}); ok {
					_, err = c.Write(buf)
					if err != nil {
						hand.log.Error(err.Error())
						c.Close()
						return
					}
				} else {
					//TODO
					hand.log.Warn("无法找到此协议。需要适配。", zap.String("RPC", string(buf)))
				}
			} else {
				hand.log.Error(err.Error())
				return
			}
		}
	}()

	return pool, nil
}

func (hand *Handle) OnMessage(
	c io.ReadWriteCloser,
	pool io.ReadWriteCloser,
	fee *fee.Fee,
	data []byte,
	id *string,
) (out []byte, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	hand.log.Info(string(data))
	req, err := ethpack.EthStratumReq(data)
	if err != nil {
		hand.log.Error(err.Error())
		c.Close()
		return
	}

	switch req.Method {
	case "mining.hello":
		fallthrough
	case "mining.subscribe":
		pool.Write(data)
		out, err = ethpack.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		return
	case "mining.authorize":
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		var worker string
		var wallet string
		//TODO
		params_zero := strings.Split(params[0], ".")
		wallet = params_zero[0]
		if len(params_zero) > 1 {
			worker = params_zero[1]
		}

		hand.log.Info("登陆矿工.", zap.String("Worker", worker), zap.String("Wallet", wallet))

		out, err = ethpack.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		pool.Write(data)
		return
	case "mining.submit":
		hand.log.Info("得到份额", zap.String("RPC", string(data)))
		pool.Write(data)

		out, err = ethpack.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		return
	default:
		hand.log.Info("KnownRpc", zap.String("RPC", string(data)))
		return
	}
}

func (hand *Handle) OnClose(id *string) {
	hand.log.Info("OnClose !!!!!")
}

func (hand *Handle) SetLog(log *zap.Logger) {
	hand.log = log
}
