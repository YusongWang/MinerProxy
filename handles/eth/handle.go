package eth

import (
	"miner_proxy/fee"
	"miner_proxy/pack"
	"miner_proxy/pack/eth"
	rpool "miner_proxy/pools"
	"miner_proxy/utils"
	"net"
	"strings"
	"sync"

	"bufio"

	"github.com/pquerna/ffjson/ffjson"
	"go.uber.org/zap"
)

var package_head = `{"id":40,"method":"eth_submitWork","params":`
var package_middle = `,"worker":"`
var package_end = `"}`

type Handle struct {
	log     *zap.Logger
	Devjob  *pack.Job
	Feejob  *pack.Job
	DevConn *net.Conn
	FeeConn *net.Conn
	SubFee  *chan []string
	SubDev  *chan []string
	Workers map[string]*pack.Worker
}

var job []string

func (hand *Handle) OnConnect(
	c net.Conn,
	config *utils.Config,
	fee *fee.Fee,
	addr string,
) (net.Conn, error) {
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

		log := hand.log.With(zap.String("Miner", c.RemoteAddr().String()))
		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				c.Close()
				hand.OnClose()
				return
			}

			go func() {
				//log.Info("收到服务器封包" + string(buf))
				var push eth.JSONPushMessage
				if err = ffjson.Unmarshal([]byte(buf), &push); err == nil {
					if result, ok := push.Result.(bool); ok {
						//增加份额
						if result {
							hand.Workers[c.LocalAddr().String()].AddShare()
							//log.Info("有效份额", zap.Any("RPC", string(buf)))
						} else {
							hand.Workers[c.LocalAddr().String()].AddReject()
							log.Warn("无效份额", zap.Any("RPC", string(buf)))
						}
					} else if params, ok := push.Result.([]interface{}); ok {
						if _, ok := hand.Workers[c.LocalAddr().String()]; !ok {
							return
						}
						hand.Workers[c.LocalAddr().String()].AddIndex()

						if utils.BaseOnIdxFee(hand.Workers[c.LocalAddr().String()].GetIndex(), rpool.DevFee) {
							if len(hand.Devjob.Job) > 0 {
								job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
							} else {
								return
							}
							res_chan := make(chan []byte)

							go func() {
								diff := utils.TargetHexToDiff(job[2])
								hand.Workers[c.LocalAddr().String()].SetDevDiff(utils.DivTheDiff(diff, hand.Workers[c.LocalAddr().String()].GetDevDiff()))
							}()

							go func() {
								fee.Dev[job[0]] = true
								job_str := ConcatJobTostr(job)
								res_chan <- ConcatToPushJob(job_str)
							}()

							job_byte := <-res_chan
							//hand.log.Info("发送开发者抽水任务", zap.String("rpc", string(job_byte)))
							_, err = c.Write(job_byte)
							if err != nil {
								log.Error(err.Error())
								c.Close()
								return
							}

						} else if utils.BaseOnIdxFee(hand.Workers[c.LocalAddr().String()].GetIndex(), config.Fee) {
							if len(hand.Feejob.Job) > 0 {
								job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
							} else {
								return
							}

							res_chan := make(chan []byte)

							go func() {
								diff := utils.TargetHexToDiff(job[2])
								hand.Workers[c.LocalAddr().String()].SetFeeDiff(utils.DivTheDiff(diff, hand.Workers[c.LocalAddr().String()].GetFeeDiff()))
							}()

							go func() {
								fee.Fee[job[0]] = true
								job_str := ConcatJobTostr(job)
								res_chan <- ConcatToPushJob(job_str)
							}()

							job_byte := <-res_chan
							//hand.log.Info("发送普通抽水任务", zap.String("rpc", string(job_byte)))
							_, err = c.Write(job_byte)
							if err != nil {
								log.Error(err.Error())
								c.Close()
								return
							}

						} else {
							go func() {
								job_params := utils.InterfaceToStrArray(params)
								diff := utils.TargetHexToDiff(job_params[2])
								hand.Workers[c.LocalAddr().String()].SetDiff(utils.DivTheDiff(diff, hand.Workers[c.LocalAddr().String()].GetDiff()))
								hand.log.Info("diff", zap.Any("diff", hand.Workers[c.LocalAddr().String()]))
								//								hand.log.Info("发送普通任务", zap.String("rpc", string(buf)))
							}()

							_, err = c.Write(buf)
							if err != nil {
								log.Error(err.Error())
								c.Close()
								return
							}

						}
					} else {
						log.Warn("无法找到此协议。需要适配。", zap.String("RPC", string(buf)))
					}
				} else {
					log.Error(err.Error())
					return
				}
			}()
		}
	}()

	return pool, nil
}

func (hand *Handle) OnMessage(
	c net.Conn,
	pool net.Conn,
	fee *fee.Fee,
	data []byte,
) (out []byte, err error) {
	req, err := eth.EthStratumReq(data)
	if err != nil {
		hand.log.Error(err.Error())
		c.Close()
		return
	}

	switch req.Method {
	case "eth_submitLogin":
		var params []string
		err = ffjson.Unmarshal(req.Params, &params)
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

		//c.LocalAddr().String() = wallet + "|" + worker

		hand.Workers[c.LocalAddr().String()] = pack.NewWorker(worker, wallet)
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		pool.Write(data)
		return
	case "eth_getWork":
		pool.Write(data)
		return
	case "eth_submitWork":
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

		var params []string
		err = ffjson.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			return
		}
		job_id := params[1]
		if _, ok := fee.Dev[job_id]; ok {
			hand.Workers[c.LocalAddr().String()].DevAdd()
			str := ConcatJobTostr(params)
			var builder strings.Builder
			builder.WriteString(package_head)
			builder.WriteString(str)
			builder.WriteString(package_middle)
			builder.WriteString("DEVFEE")
			builder.WriteString(package_end)
			builder.WriteByte('\n')

			json_rpc := builder.String()
			(*hand.DevConn).Write([]byte(json_rpc))

		} else if _, ok := fee.Fee[job_id]; ok {
			//hand.log.Info("得到普通抽水份额", zap.String("RPC", string(data)))
			hand.Workers[c.LocalAddr().String()].FeeAdd()
			str := ConcatJobTostr(params)
			var builder strings.Builder
			builder.WriteString(package_head)
			builder.WriteString(str)
			builder.WriteString(package_middle)
			builder.WriteString("MinerProxy")
			builder.WriteString(package_end)
			builder.WriteByte('\n')

			json_rpc := builder.String()
			(*hand.FeeConn).Write([]byte(json_rpc))
			//*hand.SubFee <- params
		} else {
			//hand.log.Info("得到份额", zap.String("RPC", string(data)))
			hand.Workers[c.LocalAddr().String()].AddShare()
			pool.Write(data)
		}

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
	hand.log.Info("矿机下线")
}

func (hand *Handle) SetLog(log *zap.Logger) {
	hand.log = log
}

var golbal_job = `{"id":0,"jsonrpc":"2.0","result":`
var golbal_jobend = `}`

func ConcatJobTostr(job []string) string {
	var builder strings.Builder
	builder.WriteString(`["`)

	job_len := len(job) - 1
	for i, j := range job {
		if i == job_len {
			builder.WriteString(j + `"]`)
			break
		}
		builder.WriteString(j + `","`)
	}

	return builder.String()
}

func ConcatToPushJob(job string) []byte {
	//inner_job := []byte(golbal_job + string(job) + golbal_jobend)
	var builder strings.Builder
	builder.WriteString(golbal_job)
	builder.WriteString(job)
	builder.WriteString(golbal_jobend)
	builder.WriteByte('\n')
	return []byte(builder.String())
}
