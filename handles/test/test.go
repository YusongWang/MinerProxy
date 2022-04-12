package test

import (
	"bufio"
	"io"
	"miner_proxy/fee"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Test struct {
	log *zap.Logger
}

func (hand *Test) OnConnect(
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

	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		//writer := bufio.NewWriter(c)
		//log := hand.log.With(zap.String()

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
	pool io.ReadWriteCloser,
	fee *fee.Fee,
	data []byte,
	id *string,
) (out []byte, err error) {
	hand.log.Info("矿机: " + string(data))
	pool.Write(data)
	// req, err := eth.EthStratumReq(data)
	// if err != nil {
	// 	hand.log.Error(err.Error())
	// 	c.Close()
	// 	return
	// }

	// switch req.Method {
	// case "eth_submitLogin":
	// 	var params []string
	// 	err = json.Unmarshal(req.Params, &params)
	// 	if err != nil {
	// 		hand.log.Error(err.Error())
	// 		c.Close()
	// 		return
	// 	}
	// 	var worker string
	// 	var wallet string
	// 	//TODO
	// 	params_zero := strings.Split(params[0], ".")
	// 	wallet = params_zero[0]
	// 	if len(params_zero) > 1 {
	// 		worker = params_zero[1]
	// 	} else {
	// 		if req.Worker != "" {
	// 			worker = req.Worker
	// 		}
	// 	}
	// 	hand.log.Info("登陆矿工.", zap.String("Worker", worker), zap.String("Wallet", wallet))
	// 	// reply, errReply := s.handleLoginRPC(cs, params, req.Worker)
	// 	// if errReply != nil {
	// 	// 	//return cs.sendTCPError(req.Id, errReply)
	// 	// 	log.Println("Loign Error -1")
	// 	// 	c.Close()
	// 	// 	return
	// 	// }
	// 	//return cs.sendTCPResult(req.Id, reply)
	// 	out, err = eth.EthSuccess(req.Id)
	// 	if err != nil {
	// 		hand.log.Error(err.Error())
	// 		c.Close()
	// 		return
	// 	}

	// 	pool.Write(data)
	// 	return
	// case "eth_getWork":
	// 	// reply, errReply := s.handleGetWorkRPC(cs)
	// 	// if errReply != nil {
	// 	// 	//return cs.sendTCPError(req.Id, errReply)
	// 	// 	log.Println("Loign Error -1")
	// 	// 	c.Close()
	// 	// 	return
	// 	// }
	// 	// rpc := &eth.JSONRpcResp{
	// 	// 	Id:      req.Id,
	// 	// 	Version: "2.0",
	// 	// 	Result:  true,
	// 	// }

	// 	// brpc, err := json.Marshal(rpc)
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// 	c.Close()
	// 	// 	return
	// 	// }
	// 	pool.Write(data)
	// 	// log.Println("Ret", brpc)
	// 	// out = append(brpc, '\n')
	// 	return
	// case "eth_submitWork":
	// 	var params []string
	// 	err = json.Unmarshal(req.Params, &params)
	// 	if err != nil {
	// 		hand.log.Error(err.Error())
	// 		return
	// 	}

	// 	out, err = eth.EthSuccess(req.Id)
	// 	if err != nil {
	// 		hand.log.Error(err.Error())
	// 		c.Close()
	// 		return
	// 	}

	// 	hand.log.Info("得到份额", zap.String("RPC", string(data)))
	// 	pool.Write(data)
	// 	return
	// case "eth_submitHashrate":
	// 	// 直接返回
	// 	out, err = eth.EthSuccess(req.Id)
	// 	if err != nil {
	// 		hand.log.Error(err.Error())
	// 		c.Close()
	// 		return
	// 	}

	// 	b := append(data, '\n')
	// 	pool.Write(b)
	// 	return
	// default:
	// 	hand.log.Info("KnownRpc")
	// 	return
	// }
	out = nil
	err = nil
	return
}

func (hand *Test) OnClose(id *string) {
	hand.log.Info("OnClose !!!!!")
}

func (hand *Test) SetLog(log *zap.Logger) {
	hand.log = log
}
