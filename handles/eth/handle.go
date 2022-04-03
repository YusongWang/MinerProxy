package eth

import (
	"encoding/json"
	"fmt"
	"miner_proxy/pack/eth"
	pack "miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"net"
	"strings"

	"bufio"

	"go.uber.org/zap"
	"math/rand"
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
		//writer := bufio.NewWriter(c)
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
					if rand.Intn(1000) <= int(config.Fee*10) {
						fmt.Println("收到普通抽水任务")
						var job []string
						hand.Feejob.Lock.RLock()
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
						}
						hand.Feejob.Lock.RUnlock()
						rpc := &eth.JSONPushMessage{
							Id:      0,
							Version: "2.0",
							Result:  job,
						}
						b, err := json.Marshal(rpc)
						if err != nil {
							hand.log.Error("无法序列化抽水任务", zap.Error(err))
						}
						b = append(b, '\n')
						_, err = c.Write(b)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
					} else if rand.Intn(1000) <= int(2*10) {
						fmt.Println("收到开发者抽水任务")
						var job []string
						hand.Devjob.Lock.RLock()
						if len(hand.Devjob.Job) > 0 {
							job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
						}
						hand.Devjob.Lock.RUnlock()
						rpc := &eth.JSONPushMessage{
							Id:      0,
							Version: "2.0",
							Result:  job,
						}
						b, err := json.Marshal(rpc)
						if err != nil {
							hand.log.Error("无法序列化抽水任务", zap.Error(err))
						}
						b = append(b, '\n')
						_, err = c.Write(b)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}

					} else {
						fmt.Println("收到普通任务")
						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
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

		pool.Write(data)
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

		pool.Write(data)
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

		job_id := params[1]
		fmt.Println(job_id)
		var devjob bool
		var feejob bool

		hand.Feejob.Lock.RLock()
		//TODO 优化这里的算法为O(1),目前O(n)
		for _, j := range hand.Feejob.Job {
			if j[0] == job_id {
				hand.log.Info("得到中转抽水份额", zap.String("RPC", string(data)))
				*hand.Feesub <- params
				feejob = true
				break
			}
		}
		hand.Feejob.Lock.RUnlock()

		hand.Devjob.Lock.RLock()
		for _, j := range hand.Devjob.Job {
			if j[0] == job_id {
				hand.log.Info("得到开发者抽水份额", zap.String("RPC", string(data)))
				*hand.Devsub <- params
				devjob = true
				break
			}
		}
		hand.Devjob.Lock.RUnlock()
		// TODO 判断任务JObID 是那个抽水线程的。发送到相应的抽水线程。
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		if !devjob && !feejob {
			hand.log.Info("得到份额", zap.String("RPC", string(data)))
			pool.Write(data)
		}

		return
	case "eth_submitHashrate":
		// 直接返回
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		pool.Write(data)
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
