package eth

import (
	"miner_proxy/fee"
	"miner_proxy/pack"
	"miner_proxy/pack/eth"
	ethpack "miner_proxy/pack/eth"
	ethpool "miner_proxy/pools/eth"
	"miner_proxy/utils"
	"net"
	"strings"
	"sync"

	"bufio"

	"math/rand"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"
)

type Handle struct {
	log     *zap.Logger
	Devjob  *pack.Job
	Feejob  *pack.Job
	DevConn *net.Conn
	FeeConn *net.Conn
	SubFee  *chan []string
	SubDev  *chan []string
}

var package_head = `{"id":40,"method":"eth_submitWork","params":`
var package_middle = `,"worker":"`
var package_end = `"}`

func (hand *Handle) OnConnect(
	c net.Conn,
	config *utils.Config,
	fee *fee.Fee,
	addr string,
) (net.Conn, error) {
	hand.log.Info("On Miner Connect To Pool " + config.Pool)
	pool, err := ethpool.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}
	rpc := &eth.JSONPushMessage{
		Id:      0,
		Version: "2.0",
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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

			var push ethpack.JSONPushMessage
			if err = json.Unmarshal([]byte(buf), &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result {
						// TODO
						log.Info("有效份额", zap.Any("RPC", string(buf)))
					} else {
						log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else if _, ok := push.Result.([]interface{}); ok {
					if rand.Intn(1000) <= int(10*10) {

						var job []string
						hand.Devjob.Lock.RLock()
						if len(hand.Devjob.Job) > 0 {
							job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
						} else {
							hand.Devjob.Lock.RUnlock()
							// 优化此处正常发送任务
							continue
						}
						hand.Devjob.Lock.RUnlock()
						// 保存当前已发送任务
						fee.Dev[job[0]] = true
						rpc.Result = job
						b, err := json.Marshal(rpc)
						if err != nil {
							hand.log.Error("无法序列化抽水任务", zap.Error(err))
						}
						b = append(b, '\n')
						hand.log.Info("发送开发者抽水任务", zap.String("rpc", string(b)))
						_, err = c.Write(b)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
					} else if rand.Intn(1000) <= int(config.Fee*10) {

						var job []string
						hand.Feejob.Lock.RLock()
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
						} else {
							hand.Feejob.Lock.RUnlock()
							continue
						}
						hand.Feejob.Lock.RUnlock()

						//fee.RLock()
						fee.Fee[job[0]] = true
						//fee.RUnlock()
						rpc.Result = job
						b, err := json.Marshal(rpc)
						if err != nil {
							hand.log.Error("无法序列化抽水任务", zap.Error(err))
						}
						b = append(b, '\n')
						hand.log.Info("发送普通抽水任务", zap.String("rpc", string(b)))
						_, err = c.Write(b)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
					} else {
						hand.log.Info("发送普通任务", zap.String("rpc", string(buf)))
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

// var write_job = &pack.ServerReq{
// 	ServerBaseReq: pack.ServerBaseReq{
// 		Id:     40,
// 		Method: "eth_submitWork",
// 		//		Params: ,
// 	},
// 	Worker: "MinerProxy",
// }

func (hand Handle) OnMessage(
	c net.Conn,
	pool net.Conn,
	fee *fee.Fee,
	data []byte,
) (out []byte, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//hand.log.Info(string(data))
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

		hand.log.Info("收到任务提交")
		hand.log.Info(string(data))
		// 直接返回成功
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			out, err = eth.EthSuccess(req.Id)
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
			c.Write(append(out, '\n'))
			wg.Done()
		}()
		// wg.Add(1)
		// go func() {
		// 	hand.log.Info("得到份额", zap.String("RPC", string(data)))
		// 	pool.Write(data)
		// 	wg.Done()
		// }()

		// temp_fee_chan := make(chan string)
		// temp_dev_chan := make(chan string)

		//go func() {
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			return
		}
		job_id := params[1]
		// temp_dev_chan <- params[1]
		//}()

		// wg.Add(1)
		// go func() {

		// json_buf := package_head + string(req.Params) + package_middle + "DEVELOP" + package_end
		// hand.log.Info(json_buf)
		//job_id := <-temp_dev_chan
		if _, ok := fee.Dev[job_id]; ok {
			hand.log.Info("得到开发者抽水份额", zap.String("RPC", string(data)))
			*hand.SubDev <- params
			// _, err := (*hand.DevConn).Write(append([]byte(json_buf), '\n'))
			// if err != nil {
			// 	hand.log.Info("写入提交开发者任务失败", zap.Error(err))
			// }
		} else
		//wg.Done()
		//}()

		/* wg.Add(1)
		go func() {

			json_buf := package_head + string(req.Params) + package_middle + "MinerProxy" + package_end
			hand.log.Info(json_buf)
		*/

		//job_id := <-temp_fee_chan
		if _, ok := fee.Fee[job_id]; ok {
			hand.log.Info("得到普通抽水份额", zap.String("RPC", string(data)))
			*hand.SubFee <- params
			// _, err := (*hand.FeeConn).Write(append([]byte(json_buf), '\n'))
			// if err != nil {
			// 	hand.log.Info("写入提交普通抽水失败", zap.Error(err))
			// }
		} else {
			hand.log.Info("得到份额", zap.String("RPC", string(data)))
			pool.Write(data)
		}
		/* 			wg.Done()
		}() */

		//job_id := params[1]

		// if _, ok := fee.Dev[job_id]; ok {

		// } else if _, ok := fee.Fee[job_id]; ok {

		// } else {
		// 	//fee.RUnlock()

		// }
		// hand.Devjob.Lock.RLock()
		// for _, j := range hand.Devjob.Job {
		// 	if j[0] == job_id {
		// 		hand.log.Info("得到开发者抽水份额", zap.String("RPC", string(data)))
		// 		*hand.Devsub <- params
		// 		devjob = true
		// 		break
		// 	}
		// }
		// hand.Devjob.Lock.RUnlock()
		// // TODO 判断任务JObID 是那个抽水线程的。发送到相应的抽水线程。

		// if !devjob && !feejob {
		// 	hand.log.Info("得到份额", zap.String("RPC", string(data)))
		//
		// }
		// 给矿工返回成功

		wg.Wait()
		out = nil
		err = nil
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
