package eth

import (
	"encoding/json"
	"fmt"
	"miner_proxy/pack/eth"
	pack "miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"net"

	"bufio"

	"go.uber.org/zap"
)

type Handle struct {
	log    *zap.Logger
	Devjob *eth.Job
	Feejob *eth.Job
	Devsub *chan []string
	Feesub *chan []string
}

func (hand *Handle) OnConnect(c net.Conn, config *utils.Config, addr string) (net.Conn, error) {
	hand.log.Info("On Miner Connect To Pool " + config.Pool)
	pool, err := ethpool.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}

	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		log := hand.log.With(zap.String("Miner", c.RemoteAddr().String()))
		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			var push pack.JSONPushMessage
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
					fmt.Println("收到普通任务")
					b := append(buf, '\n')
					_, err = c.Write(b)
					if err != nil {
						log.Error(err.Error())
						c.Close()
						return
					}
				} else {
					//TODO
				}
			} else {
				log.Error(err.Error())
				return
			}
		}
	}()

	return pool, nil
}

func (hand *Handle) OnMessage(c net.Conn, pool net.Conn, data []byte) (out []byte, err error) {
	hand.log.Info(string(data))
	req, err := eth.EthStratumReq(data)
	if err != nil {
		hand.log.Error(err.Error())
		c.Close()
		return
	}

	switch req.Method {
	case "eth_submitLogin":
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		// l := s.B.Listen()
		// go func() {
		// 	for {
		// 		select {
		// 		case job := <-l.Ch:
		// 			//fmt.Println(job)
		// 			c.Send(job)
		// 		}
		// 	}
		// }()

		// if req.Worker != "" {
		// 	s.Worker = req.Worker
		// } else {
		// 	p1 := strings.Split(params[0], ".")
		// }

		// reply, errReply := s.handleLoginRPC(cs, params, req.Worker)
		// if errReply != nil {
		// 	//return cs.sendTCPError(req.Id, errReply)
		// 	log.Println("Loign Error -1")
		// 	c.Close()
		// 	return
		// }
		//return cs.sendTCPResult(req.Id, reply)
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		b := append(data, '\n')
		pool.Write(b)
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
		b := append(data, '\n')
		pool.Write(b)

		// log.Println("Ret", brpc)
		// out = append(brpc, '\n')
		return
	case "eth_submitWork":
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			return
		}
		//s.Remote <- params
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		b := append(data, '\n')
		pool.Write(b)

		return
	case "eth_submitHashrate":
		// 直接返回
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		b := append(data, '\n')
		pool.Write(b)
		return
	default:
		hand.log.Info("KnownRpc")
		return
	}
}

func (hand *Handle) OnClose() {
	hand.log.Info("OnClose !!!!!")
}

func (hand *Handle) SetLog(log *zap.Logger) {
	hand.log = log
}
