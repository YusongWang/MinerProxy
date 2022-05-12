package eth

import (
	"io"

	"miner_proxy/global"
	pack "miner_proxy/pack"
	"miner_proxy/pack/eth"
	"miner_proxy/utils"
	"strings"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"

	"bufio"

	"go.uber.org/zap"
)

type NoFeeHandle struct {
	log *zap.Logger
}

func (hand *NoFeeHandle) OnConnect(
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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		//writer := bufio.NewWriter(c)
		log := hand.log

		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			var push eth.JSONPushMessage
			if err = json.Unmarshal([]byte(buf), &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result == true {
						// TODO
						log.Info("有效份额", zap.Any("RPC", string(buf)))
					} else {
						log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else if _, ok := push.Result.([]interface{}); ok {

					b := append(buf, '\n')
					_, err = c.Write(b)
					if err != nil {
						log.Error(err.Error())
						c.Close()
						return
					}

				} else {
					//TODO
					log.Warn("无法找到此协议。需要适配。", zap.String("RPC", string(buf)))
				}
			} else {
				log.Error(err.Error())
				return
			}
		}
	}()

	return pool, nil
}

func (hand *NoFeeHandle) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
	data *[]byte,
	worker *pack.Worker,
) (out []byte, err error) {
	hand.log.Info(string(*data))
	req, err := eth.EthStratumReq(*data)
	if err != nil {
		hand.log.Error(err.Error())
		c.Close()
		return
	}

	switch req.Method {
	case "eth_submitLogin":
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
		} else {
			if req.Worker != "" {
				worker = req.Worker
			}
		}
		hand.log.Info("登陆矿工.", zap.String("Worker", worker), zap.String("Wallet", wallet))

		var id int64
		id, err = jsonparser.GetInt(*data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		out, err = eth.EthSuccess(id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		(*pool).Write(*data)
		return
	case "eth_getWork":
		// reply, errReply := s.handleGetWorkRPC(cs)
		// if errReply != nil {
		// 	//return cs.sendTCPError(req.Id, errReply)
		// 	log.Println("Loign Error -1")
		// 	c.Close()
		// 	return
		// }
		// rpc := &eth.JSONRpcResp{
		// 	Id:      req.Id,
		// 	Version: "2.0",
		// 	Result:  true,
		// }

		// brpc, err := json.Marshal(rpc)
		// if err != nil {
		// 	log.Println(err)
		// 	c.Close()
		// 	return
		// }
		(*pool).Write(*data)
		// log.Println("Ret", brpc)
		// out = append(brpc, '\n')
		return
	case "eth_submitWork":
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			return
		}
		var id int64
		id, err = jsonparser.GetInt(*data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		out, err = eth.EthSuccess(id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		hand.log.Info("得到份额", zap.String("RPC", string(*data)))
		(*pool).Write(*data)
		return
	case "eth_submitHashrate":
		var id int64
		id, err = jsonparser.GetInt(*data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		// 直接返回
		out, err = eth.EthSuccess(id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		b := append(*data, '\n')
		(*pool).Write(b)
		return
	default:
		hand.log.Info("KnownRpc")
		return
	}
}

func (hand *NoFeeHandle) OnClose(worker *pack.Worker) {
	hand.log.Info("OnClose !!!!!")
}

func (hand *NoFeeHandle) SetLog(log *zap.Logger) {
	hand.log = log
}
