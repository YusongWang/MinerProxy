package eth

import (
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

	// var rpc_id int64
	// rpc_id, err = jsonparser.GetInt(*data, "id")
	// if err != nil {
	// 	err = nil
	// 	rpc_id = 0
	// }

	switch method {
	case "mining.subscribe":
		worker.SetProtocol(eth.ProtocolLegacyStratum)

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

		// 读取Miner内核客户端
		if len(params) >= 1 {
			worker.SetClientAgent(params[0])
		}

		if len(params) >= 2 {
			if strings.HasPrefix(strings.ToLower(params[1]), eth.EthereumStratumPrefix) {
				worker.SetProtocol(eth.ProtocolEthereumStratum)
			}
		}

		//TODO 解析协议，
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
		} else {
			worker.AddShare()
			_, err = (*pool).Write(*data)
			if err != nil {
				hand.log.Error("写入矿池失败: " + err.Error())
				c.Close()
				return
			}
		}

		c.Write(out)
		out = nil
		err = nil
		return
	case "mining.extranonce.subscribe":
		fallthrough
	case "eth_submitHashrate":
		{
			hashrate, innerErr := jsonparser.GetString(*data, "params", "[0]")
			if innerErr == nil {
				worker.SetReportHash(utils.String2Big(hashrate))
			}
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
