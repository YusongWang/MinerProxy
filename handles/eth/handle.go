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
	Workers map[string]*pack.Worker
	Wallet  string
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
			log.Info("收到服务器封包" + string(buf))
			var push eth.JSONPushMessage
			if err = json.Unmarshal([]byte(buf), &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result {
						log.Info("有效份额", zap.Any("RPC", string(buf)))
					} else {
						log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else if params, ok := push.Result.([]interface{}); ok {
					hand.Workers[hand.Wallet].AddIndex()
					if utils.BaseOnIdxFee(hand.Workers[hand.Wallet].GetIndex(), rpool.DevFee) {
						if len(hand.Devjob.Job) > 0 {
							job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
						} else {
							continue
						}
						fee.Dev[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)
						hand.log.Info("发送开发者抽水任务", zap.String("rpc", string(job_byte)))
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
					} else if utils.BaseOnIdxFee(hand.Workers[hand.Wallet].GetIndex(), config.Fee) {
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
						} else {
							continue
						}
						fee.Fee[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)
						hand.log.Info("发送普通抽水任务", zap.String("rpc", string(job_byte)))
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							return
						}
					} else {
						job_params := utils.InterfaceToStrArray(params)
						diff := utils.TargetHexToDiff(job_params[2])
						hand.Workers[hand.Wallet].SetDiff(utils.DivTheDiff(diff, hand.Workers[hand.Wallet].GetDiff()))
						hand.log.Info("diff", zap.Any("diff", hand.Workers[hand.Wallet]))

						hand.log.Info("发送普通任务", zap.String("rpc", string(buf)))
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
		}
	}()

	return pool, nil
}

func (hand Handle) OnMessage(
	c net.Conn,
	pool net.Conn,
	fee *fee.Fee,
	data []byte,
) (out []byte, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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

		hand.Wallet = wallet
		hand.Workers[wallet] = pack.NewWorker(worker, wallet)

		hand.log.Info("登陆矿工.", zap.String("Worker", worker), zap.String("Wallet", wallet))

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

		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			hand.log.Error(err.Error())
			return
		}
		job_id := params[1]
		if _, ok := fee.Dev[job_id]; ok {
			hand.log.Info("得到开发者抽水份额", zap.String("RPC", string(data)))
			hand.Workers[hand.Wallet].DevAdd()
			*hand.SubDev <- params
		} else if _, ok := fee.Fee[job_id]; ok {
			hand.log.Info("得到普通抽水份额", zap.String("RPC", string(data)))
			hand.Workers[hand.Wallet].FeeAdd()
			//(*hand.Feejob).Write()
			*hand.SubFee <- params
		} else {
			hand.log.Info("得到份额", zap.String("RPC", string(data)))
			hand.Workers[hand.Wallet].AddShare()
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
	hand.log.Info("OnClose !!!!!")
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
