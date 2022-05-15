package eth

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"miner_proxy/global"
	"miner_proxy/pack/eth"
	pools "miner_proxy/pools"
	"miner_proxy/utils"
	"strings"

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
	Devjob  *global.Job
	Feejob  *global.Job
	DevConn *io.ReadWriteCloser
	FeeConn *io.ReadWriteCloser
	SubFee  *chan []byte
	SubDev  *chan []byte
}

var job []string

func (hand *Handle) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
	addr string,
	worker *global.Worker,
) (io.ReadWriteCloser, error) {
	return nil, nil
}

func (hand *Handle) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	proxyFee *global.Fee,
	data *[]byte,
	worker *global.Worker,
) (out []byte, err error) {
	method, err := jsonparser.GetString(*data, "method")
	if err != nil {
		hand.log.Info("非法封包", zap.String("package", string(*data)))
		c.Close()
		return
	}

	var rpc_id int64
	rpc_id, err = jsonparser.GetInt(*data, "id")
	if err != nil {
		err = nil
		rpc_id = 0
	}

	switch method {
	case "mining.subscribe":
		// 直接返回
		out, err = eth.EthSuccess(rpc_id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		//TODO 解析协议，
		//TODO 解析客户端 miner
		*pool, err = ConnectToPool(c, hand, config, proxyFee, worker)
		if err != nil {
			hand.log.Error("矿池拒绝链接或矿池地址不正确! " + err.Error())
			return
		}
		worker.SetAuthStat(eth.StatSubScribed)
		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}

		return
	case "mining.authorize":
		fallthrough
	case "eth_submitLogin":
		var params []string
		var parse_byte []byte
		var name string
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

		name, _ = jsonparser.GetString(*data, "worker")

		if worker.AuthorizeStat == eth.StatConnected {
			*pool, err = ConnectToPool(c, hand, config, proxyFee, worker)
			if err != nil {
				hand.log.Error("矿池拒绝链接或矿池地址不正确! " + err.Error())
				return
			}
			worker.SetAuthStat(eth.StatSubScribed)
		}

		if !worker.Authorize(method, params, name) {
			err = errors.New("矿工登录失败")
			return
		}

		global.GonlineWorkers.Lock()
		global.GonlineWorkers.Workers[worker.Fullname] = worker
		global.GonlineWorkers.Unlock()

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
	case "mining.set_difficulty":
		fallthrough
	case "eth_getWork":
		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}
		return
	case "mining.submit":
		fallthrough
	case "eth_submitWork":

		var job_id string
		job_id, err = jsonparser.GetString(*data, "params", "[1]")
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		if _, ok := proxyFee.Dev.Load(job_id); ok {
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
		} else if _, ok := proxyFee.Fee.Load(job_id); ok {
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

		out, err = eth.EthSuccess(rpc_id)
		if err != nil {
			hand.log.Error(err.Error())
			c.Close()
			return
		}

		c.Write(out)
		out = nil
		err = nil
		return
	case "mining.extranonce.subscribe":
		fallthrough
	case "eth_submitHashrate":
		{
			var hashrate string
			hashrate, err = jsonparser.GetString(*data, "params", "[0]")
			if err == nil {
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
	// ignore unimplemented methods
	case "mining.multi_version":
		fallthrough
	case "mining.suggest_difficulty":
		// If no response, the miner may wait indefinitely
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

func (hand *Handle) OnClose(worker *global.Worker) {
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
	proxyFee *global.Fee,
	worker *global.Worker,
) (pool io.ReadWriteCloser, err error) {
	pool, err = utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}

	log := (*hand.log).With(zap.String("UUID", worker.Id), zap.String("wallet", worker.Wallet), zap.String("worker", worker.Worker_name))

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

						proxyFee.Dev.Store(job[0], global.FeeResult{})

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

						proxyFee.Fee.Store(job[0], global.FeeResult{})

						job_str := ConcatJobTostr(job)
						job_byte := ConcatToPushJob(job_str)

						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							pool.Close()

							return
						}

					} else {

						job_diff, err := jsonparser.GetString(buf, "result", "[2]")
						if err != nil {
							log.Info("格式化Diff字段失败")
							log.Error(err.Error())
							c.Close()
							pool.Close()
							return
						}

						diff := utils.TargetHexToDiff(job_diff)
						worker.SetDiff(diff)

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
