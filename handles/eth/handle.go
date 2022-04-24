package eth

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"miner_proxy/fee"
	"miner_proxy/pack"
	"miner_proxy/pack/eth"
	pools "miner_proxy/pools"
	"miner_proxy/utils"
	"strings"
	"sync"

	"github.com/buger/jsonparser"
	"go.uber.org/zap"
)

// var package_head = `{"id":40,"method":"eth_submitWork","params":`
// var package_middle = `,"worker":"`
// var package_end = `"}`

type Handle struct {
	log     *zap.Logger
	Devjob  *pack.Job
	Feejob  *pack.Job
	DevConn *io.ReadWriteCloser
	FeeConn *io.ReadWriteCloser
	SubFee  *chan []byte
	SubDev  *chan []byte
	Workers map[string]*pack.Worker
}

var job []string

func (hand *Handle) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *fee.Fee,
	addr string,
	id *string,
) (io.ReadWriteCloser, error) {
	return nil, nil
}

func (hand *Handle) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	fee *fee.Fee,
	data []byte,
	id *string,
) (out []byte, err error) {
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	defer func() {
		if x := recover(); x != nil {
			hand.log.Info("Recover", zap.Any("err", x))
			err = errors.New("Panic(). first package not the Login")
			return
		}
	}()

	method, err := jsonparser.GetString(data, "method")
	if err != nil {
		hand.log.Info("非法封包", zap.String("package", string(data)))
		c.Close()
		return
	}

	switch method {
	case "eth_submitLogin":
		var params []string
		var parse_byte []byte
		parse_byte, _, _, err = jsonparser.Get(data, "params")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		err = json.Unmarshal(parse_byte, &params)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		var worker string
		var wallet string

		params_zero := strings.Split(params[0], ".")
		wallet = params_zero[0]
		if len(params_zero) > 1 {
			worker = params_zero[1]
		} else {
			worker, err = jsonparser.GetString(data, "worker")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
		}
		var pool_conn io.ReadWriteCloser
		pool_conn, err = ConnectToPool(c, hand, config, fee, id)
		if err != nil {
			hand.log.Error("矿池拒绝链接或矿池地址不正确! " + err.Error())
			return
		}

		*pool = pool_conn
		hand.Workers[*id] = pack.NewWorker(worker, wallet, *id)
		hand.Workers[*id].Logind(worker, wallet)

		hand.log.Info("登陆矿工.", zap.String("Worker", worker), zap.String("Wallet", wallet))

		var rpc_id int64
		rpc_id, err = jsonparser.GetInt(data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		out, err = eth.EthSuccess(rpc_id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		_, err = (*pool).Write(data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			hand.OnClose(id)
			return
		}

		return
	case "eth_getWork":
		if _, ok := hand.Workers[*id]; !ok {
			return
		}
		_, err = (*pool).Write(data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			hand.OnClose(id)
			return
		}
		return
	case "eth_submitWork":
		if _, ok := hand.Workers[*id]; !ok {
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			id, err := jsonparser.GetInt(data, "id")
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
			c.Write(out)
			wg.Done()
		}()

		go func() {
			var job_id string
			job_id, err = jsonparser.GetString(data, "params", "[1]")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}

			if _, ok := fee.Dev[job_id]; ok {
				hand.Workers[*id].DevAdd()

				var parse_byte []byte
				parse_byte, _, _, err = jsonparser.Get(data, "params")
				if err != nil {
					hand.log.Error(err.Error())
					c.Close()
					return
				}
				*hand.SubDev <- parse_byte
			} else if _, ok := fee.Fee[job_id]; ok {

				hand.Workers[*id].FeeAdd()
				var parse_byte []byte
				parse_byte, _, _, err = jsonparser.Get(data, "params")
				if err != nil {
					hand.log.Error(err.Error())
					c.Close()
					return
				}
				*hand.SubFee <- parse_byte
			} else {
				hand.Workers[*id].AddShare()
				_, err = (*pool).Write(data)
				if err != nil {
					hand.log.Error("写入矿池失败: " + err.Error())
					c.Close()
					hand.OnClose(id)
					return
				}
			}
		}()

		wg.Wait()
		out = nil
		err = nil
		return
	case "eth_submitHashrate":
		if _, ok := hand.Workers[*id]; !ok {
			return
		}
		var rpc_id int64
		rpc_id, err = jsonparser.GetInt(data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		{
			var hashrate string
			hashrate, err = jsonparser.GetString(data, "params", "[0]")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
			hand.Workers[*id].SetReportHash(utils.String2Big(hashrate))
		}
		// 直接返回
		out, err = eth.EthSuccess(rpc_id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		_, err = (*pool).Write(data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			hand.OnClose(id)
			return
		}
		return
	default:
		hand.log.Info("KnownRpc")
		return
	}
}

func (hand *Handle) OnClose(id *string) {
	if worker, ok := hand.Workers[*id]; ok {
		if worker.IsOnline {
			worker.Logout()
			hand.log.Info("矿机下线", zap.Any("Worker", worker))
		}
	}
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

func ConnectToPool(
	c io.ReadWriteCloser,
	hand *Handle,
	config *utils.Config,
	fee *fee.Fee,
	id *string,
) (pool io.ReadWriteCloser, err error) {
	defer func() {
		if x := recover(); x != nil {
			hand.log.Info("Recover", zap.Any("err", x))
			c.Close()
			err = errors.New("Panic() Recover. ConnectToPool To Pool Error")
			return
		}
	}()

	pool, err = utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}

	log := (*hand.log).With(zap.String("UUID", *id))

	reader := bufio.NewReader(pool)
	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func(reader *bufio.Reader) {
		defer func() {
			if x := recover(); x != nil {
				hand.log.Info("Recover", zap.Any("err", x))
				c.Close()
				err = errors.New("Panic() Race!!!!!!! . But Why ????")
				return
			}
		}()

		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				c.Close()
				hand.OnClose(id)
				return
			}

			if result, _, _, err := jsonparser.Get(buf, "result"); err == nil {
				//if result, ok := buf.(bool); ok {
				if res, err := jsonparser.ParseBoolean(result); err == nil {
					//增加份额
					if res {
						if w, ok := hand.Workers[*id]; ok {
							w.AddShare()
						}
						//log.Info("有效份额", zap.Any("RPC", string(buf)))
					} else {
						if w, ok := hand.Workers[*id]; ok {
							w.AddReject()
						}
						log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else {
					if _, ok := hand.Workers[*id]; !ok {
						continue
					}
					// if w, ok := hand.Workers[*id]; ok {
					// 	w.AddShare()
					// }
					hand.Workers[*id].AddIndex()
					if utils.BaseOnIdxFee(hand.Workers[*id].GetIndex(), pools.DevFee) {
						if len(hand.Devjob.Job) > 0 {
							job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
						} else {
							continue
						}

						if len(job) == 0 {
							log.Info("当前job内容为空")
							continue
						}
						diff := utils.TargetHexToDiff(job[2])
						hand.Workers[*id].SetDevDiff(utils.DivTheDiff(diff, hand.Workers[*id].GetDevDiff()))

						fee.Dev[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)

						//job_byte := <-res_chan
						//log.Info("发送开发者抽水任务", zap.String("rpc", string(job_byte)))
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							hand.OnClose(id)
							c.Close()
							return
						}

					} else if utils.BaseOnIdxFee(hand.Workers[*id].GetIndex(), config.Fee) {
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
							//log.Info("得到当前Job", zap.Any("job", job))
						} else {
							continue
						}

						if len(job) == 0 {
							log.Info("当前job内容为空")
							continue
						}
						diff := utils.TargetHexToDiff(job[2])
						hand.Workers[*id].SetFeeDiff(utils.DivTheDiff(diff, hand.Workers[*id].GetFeeDiff()))

						fee.Fee[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)

						//log.Info("发送普通抽水任务", zap.String("rpc", string(job_byte)))
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							hand.OnClose(id)
							return
						}

					} else {
						//go func() {
						job_diff, err := jsonparser.GetString(buf, "result", "[2]")
						if err != nil {
							log.Info("格式化Diff字段失败")
							log.Error(err.Error())
							hand.OnClose(id)
							c.Close()
							return
						}

						diff := utils.TargetHexToDiff(job_diff)
						hand.Workers[*id].SetDiff(utils.DivTheDiff(diff, hand.Workers[*id].GetDiff()))
						// log.Info("diff", zap.Any("diff", hand.Workers[*id]))
						// log.Info("发送普通任务", zap.String("rpc", string(buf)))
						//}()
						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())
							hand.OnClose(id)
							c.Close()
							return
						}

					}
				}
			} else {
				c.Close()
				hand.OnClose(id)
				log.Error(err.Error())
				return
			}
		}
	}(reader)

	return
}
