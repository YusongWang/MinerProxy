package eth

import (
	"encoding/json"
	"errors"
	"io"

	"miner_proxy/global"
	"miner_proxy/pack/eth"
	"miner_proxy/utils"
	"strings"

	"github.com/buger/jsonparser"

	"bufio"

	"go.uber.org/zap"
)

type NoFeeHandle struct {
	log *zap.Logger
}

func (hand *NoFeeHandle) OnConnect(
	c io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
	addr string,
	worker *global.Worker,
) (io.ReadWriteCloser, error) {
	hand.log.Info("On Miner Connect To Pool " + config.Pool)
	pool, err := utils.NewPool(config.Pool)
	if err != nil {
		hand.log.Warn("矿池连接失败", zap.Error(err), zap.String("pool", config.Pool))
		c.Close()
		return nil, err
	}
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// 处理上游矿池。如果连接失败。矿工线程直接退出并关闭
	go func() {
		reader := bufio.NewReader(pool)
		//writer := bufio.NewWriter(c)
		log := hand.log

		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			// Share

			b := append(buf, '\n')
			_, err = c.Write(b)
			if err != nil {
				log.Error(err.Error())
				c.Close()
				return
			}
		}
	}()

	return pool, nil
}

func (hand *NoFeeHandle) OnMessage(
	c io.ReadWriteCloser,
	pool *io.ReadWriteCloser,
	config *utils.Config,
	fee *global.Fee,
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
		*pool, err = ConnectNoFeePool(c, hand, config, worker)
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
	case "eth_submitLogin":
		worker.SetAuthStat(eth.StatSubScribed)
		worker.SetProtocol(eth.ProtocolETHProxy)
		fallthrough
	case "mining.authorize":
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

		if worker.AuthorizeStat == eth.StatSubScribed && worker.Protocol == eth.ProtocolETHProxy {
			*pool, err = ConnectNoFeePool(c, hand, config, worker)
			if err != nil {
				hand.log.Error("矿池拒绝链接或矿池地址不正确! " + err.Error())
				return
			}
		}

		if !worker.Authorize(method, params, name) {
			err = errors.New("矿工登录失败")
			return
		}
		worker.SetAuthStat(eth.StatAuthorized)

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

		//utils.Logger.Info("Submit", zap.Int("protocol", int(worker.Protocol)), zap.String("powHash", powHash), zap.String("mixHash", mixHash), zap.String("nonce", nonce))
		// var job_id string
		// job_id, err = jsonparser.GetString(*data, "params", "[1]")
		// if err != nil {
		// 	hand.log.Error(err.Error())
		// 	c.Close()
		// 	return
		// }

		worker.AddShare()
		_, err = (*pool).Write(*data)
		if err != nil {
			hand.log.Error("写入矿池失败: " + err.Error())
			c.Close()
			return
		}

		//c.Write(out)
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

func (hand *NoFeeHandle) OnClose(worker *global.Worker) {
	hand.log.Info("OnClose !!!!!")
}

func (hand *NoFeeHandle) SetLog(log *zap.Logger) {
	hand.log = log
}

func ConnectNoFeePool(
	c io.ReadWriteCloser,
	hand *NoFeeHandle,
	config *utils.Config,
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

			//log.Info("Message", zap.String("RPC", string(buf)))

			if worker.Protocol == eth.ProtocolETHProxy {
				if result, _, _, err := jsonparser.Get(buf, "result"); err == nil {
					//if result, ok := buf.(bool); ok {
					if res, err := jsonparser.ParseBoolean(result); err == nil {
						//增加份额
						if res {
							worker.AddShare()
						} else {
							worker.AddReject()
						}

						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())

							c.Close()
							pool.Close()
							return
						}
					} else {
						worker.AddIndex()
						if worker.Worker_idx == 5 {
							job_diff, err := jsonparser.GetString(buf, "result", "[2]")
							if err == nil {
								diff := utils.TargetHexToDiff(job_diff)
								worker.SetDiff(diff)
							}
						}
						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())

							c.Close()
							pool.Close()
							return
						}
					}
				} else {
					c.Close()
					pool.Close()
					log.Error(err.Error())
					return
				}
			} else if worker.Protocol == eth.ProtocolLegacyStratum {
				if result, _, _, err := jsonparser.Get(buf, "result"); err == nil {
					if res, err := jsonparser.ParseBoolean(result); err == nil {
						// 增加份额
						if res {
							worker.AddShare()
						} else {
							worker.AddReject()
						}

						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())

							c.Close()
							pool.Close()
							return
						}
					}
				} else if _, _, _, err := jsonparser.Get(buf, "params"); err == nil {
					worker.AddIndex()

					if worker.Worker_idx == 5 {
						job_diff, err := jsonparser.GetString(buf, "params", "[2]")
						if err == nil {
							diff := utils.TargetHexToDiff(job_diff)
							worker.SetDiff(diff)
						}
					}
					_, err = c.Write(buf)
					if err != nil {
						log.Error(err.Error())

						c.Close()
						pool.Close()
						return
					}

				} else {
					c.Close()
					pool.Close()
					log.Error(err.Error())
					return
				}
			} else if worker.Protocol == eth.ProtocolEthereumStratum {
				if result, _, _, err := jsonparser.Get(buf, "result"); err == nil {
					if res, err := jsonparser.ParseBoolean(result); err == nil {
						// 增加份额
						if res {
							worker.AddShare()
						} else {
							worker.AddReject()
						}

						_, err = c.Write(buf)
						if err != nil {
							log.Error(err.Error())

							c.Close()
							pool.Close()
							return
						}
					}
				} else if _, _, _, err := jsonparser.Get(buf, "params"); err == nil {
					worker.AddIndex()

					if worker.Worker_idx == 5 {
						job_diff, err := jsonparser.GetString(buf, "params", "[2]")
						if err == nil {
							diff := utils.TargetHexToDiff(job_diff)
							worker.SetDiff(diff)
						}
					}
					_, err = c.Write(buf)
					if err != nil {
						log.Error(err.Error())

						c.Close()
						pool.Close()
						return
					}

				} else {
					c.Close()
					pool.Close()
					log.Error(err.Error())
					return
				}
			} else {
				_, err = c.Write(buf)
				if err != nil {
					log.Error(err.Error())

					c.Close()
					pool.Close()
					return
				}
			}
		}
	}(reader)

	return
}
