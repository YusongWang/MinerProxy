package eth

import (
	"bufio"
	"encoding/json"
	"io"
	"miner_proxy/fee"
	"miner_proxy/global"
	"miner_proxy/pack"
	"miner_proxy/pack/eth"
	pools "miner_proxy/pools"
	"miner_proxy/utils"
	"strings"
	"sync"

	//"github.com/pkg/profile"

	"github.com/buger/jsonparser"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

var package_head = `{"id":40,"method":"eth_submitWork","params":`
var package_middle = `,"worker":"`
var package_end = `"}`

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
}

var job []string

func (hand *Handle) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *fee.Fee,
	addr string,
	worker *pack.Worker,
) (io.ReadWriteCloser, error) {
	return nil, nil
}

func (hand *Handle) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	fee *fee.Fee,
	data *[]byte,
	worker *pack.Worker,
) (out []byte, err error) {
	method, err := jsonparser.GetString(*data, "method")
	if err != nil {
		hand.log.Info("非法封包", zap.String("package", string(*data)))
		c.Close()
		return
	}

	switch method {
	case "eth_submitLogin":
		var params []string
		var parse_byte []byte
		parse_byte, _, _, err = jsonparser.Get(*data, "params")
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

		var worker_name string
		var wallet string

		params_zero := strings.Split(params[0], ".")
		wallet = params_zero[0]
		if len(params_zero) > 1 {
			worker_name = params_zero[1]
		} else {
			worker_name, err = jsonparser.GetString(*data, "worker")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
		}

		*pool, err = ConnectToPool(c, hand, config, fee, worker, wallet, worker_name)
		if err != nil {
			hand.log.Error("矿池拒绝链接或矿池地址不正确! " + err.Error())
			return
		}
		if worker_name == "" {
			worker_name = "default"
		}

		worker_info := wallet + "." + worker_name

		worker.Logind(worker_name, wallet)

		global.GonlineWorkers.Lock()
		global.GonlineWorkers.Workers[worker_info] = worker
		global.GonlineWorkers.Unlock()

		//hand.log.Info("登陆矿工.", zap.String("Worker", worker_name), zap.String("Wallet", wallet))

		var rpc_id int64
		rpc_id, err = jsonparser.GetInt(*data, "id")
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

		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}

		return
	case "eth_getWork":
		// if _, ok := hand.Workers[*id]; !ok {
		// 	return
		// }
		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}
		return
	case "eth_submitWork":
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			id, err := jsonparser.GetInt(*data, "id")
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

		var job_id string
		job_id, err = jsonparser.GetString(*data, "params", "[1]")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		if _, ok := fee.Dev[job_id]; ok {
			worker.DevAdd()
			var parse_byte []byte
			parse_byte, _, _, err = jsonparser.Get(*data, "params")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
			var builder strings.Builder
			builder.WriteString(package_head)
			builder.Write(parse_byte)
			builder.WriteString(package_middle)
			builder.WriteString(pools.DEVELOP)
			builder.WriteString(package_end)
			builder.WriteByte('\n')

			json_rpc := builder.String()
			_, err = (*hand.DevConn).Write([]byte(json_rpc))
			if err != nil {
				hand.log.Error("写入矿池失败: " + err.Error())
				c.Close()
				return
			}
			//(*hand.DevConn).Flush()
		} else if _, ok := fee.Fee[job_id]; ok {
			worker.FeeAdd()
			var parse_byte []byte
			parse_byte, _, _, err = jsonparser.Get(*data, "params")
			if err != nil {
				hand.log.Error(err.Error())
				c.Close()
				return
			}
			var builder strings.Builder
			builder.WriteString(package_head)
			builder.Write(parse_byte)
			builder.WriteString(package_middle)
			builder.WriteString(config.Worker)
			builder.WriteString(package_end)
			builder.WriteByte('\n')
			json_rpc := builder.String()
			_, err = (*hand.FeeConn).Write([]byte(json_rpc))
			if err != nil {
				hand.log.Error("写入矿池失败: " + err.Error())
				c.Close()
				return
			}
			//(*hand.FeeConn).Flush()
			//*hand.SubFee <- parse_byte
		} else {
			worker.AddShare()
			_, err = (*pool).Write(*data)
			if err != nil {
				hand.log.Error("写入矿池失败: " + err.Error())
				c.Close()
				return
			}
		}

		wg.Wait()
		out = nil
		err = nil
		return
	case "eth_submitHashrate":
		var rpc_id int64
		rpc_id, err = jsonparser.GetInt(*data, "id")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		{
			var hashrate string
			hashrate, err = jsonparser.GetString(*data, "params", "[0]")
			if err != nil {
				hand.log.Error(err.Error())
			} else {
				worker.SetReportHash(utils.String2Big(hashrate))
			}
		}

		// 直接返回
		out, err = eth.EthSuccess(rpc_id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}
		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}
		return
	default:
		hand.log.Info("KnownRpc")
		return
	}
}

func (hand *Handle) OnClose(worker *pack.Worker) {
	if worker.IsOnline() {
		worker.Logout()
		hand.log.Info("矿机下线", zap.Any("Worker", worker), zap.String("Time", humanize.Time(worker.Login_time)))
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
	worker *pack.Worker,
	wallet string,
	worker_name string,
) (pool io.ReadWriteCloser, err error) {
	pool, err = utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}
	log := (*hand.log).With(zap.String("UUID", worker.Id), zap.String("wallet", wallet), zap.String("worker", worker_name))

	reader := bufio.NewReader(pool)
	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func(read *bufio.Reader) {
		var buf []byte
		for {
			buf, err = read.ReadBytes('\n')
			if err != nil {
				c.Close()
				pool.Close()
				return
			}

			if result, _, _, err := jsonparser.Get(buf, "result"); err == nil {
				//if result, ok := buf.(bool); ok {
				if res, err := jsonparser.ParseBoolean(result); err == nil {
					//增加份额
					if res {
						worker.AddShare()
					} else {
						worker.AddReject()
						log.Warn("无效份额", zap.Any("RPC", string(buf)))
					}
				} else {
					worker.AddIndex()

					if utils.BaseOnRandFee(worker.GetIndex(), pools.DevFee) {
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
						worker.SetDevDiff(diff)

						fee.Dev[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)

						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							pool.Close()
							return
						}

					} else if utils.BaseOnRandFee(worker.GetIndex(), config.Fee) {
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
						} else {
							continue
						}

						if len(job) == 0 {
							log.Info("当前job内容为空")
							continue
						}
						diff := utils.TargetHexToDiff(job[2])
						worker.SetFeeDiff(diff)

						fee.Fee[job[0]] = true
						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)

						//log.Info("发送普通抽水任务", zap.String("rpc", string(job_byte)))
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							pool.Close()

							return
						}

					} else {
						//go func() {
						job_diff, err := jsonparser.GetString(buf, "result", "[2]")
						if err != nil {
							log.Info("格式化Diff字段失败")
							log.Error(err.Error())
							c.Close()
							pool.Close()
							return
						}

						diff := utils.TargetHexToDiff(job_diff)
						//worker.SetDiff(utils.DivTheDiff(diff, worker.GetDiff()))
						worker.SetDiff(diff)
						// log.Info("diff", zap.Any("diff", worker))
						// log.Info("发送普通任务", zap.String("rpc", string(buf)))
						//}()
						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())

							c.Close()
							pool.Close()
							return
						}

					}
				}
			} else {
				c.Close()
				pool.Close()
				log.Error(err.Error())
				return
			}
		}
	}(reader)

	return
}
