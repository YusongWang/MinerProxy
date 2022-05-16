package eth

import (
	"bufio"
	"encoding/json"
	"io"
	"miner_proxy/global"
	"miner_proxy/pack/eth"
	pools "miner_proxy/pools"
	"miner_proxy/utils"
	"strings"

	"github.com/buger/jsonparser"
	"go.uber.org/zap"
)

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

			log.Info("Message", zap.String("RPC", string(buf)))

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
						if utils.BaseOnRandFee(worker.GetIndex(), pools.DevFee) {
							if len(hand.Devjob.Job) > 0 {
								job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
							} else {
								goto SendWorker
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
							continue
						} else if utils.BaseOnRandFee(worker.GetIndex(), config.Fee) {
							if len(hand.Feejob.Job) > 0 {
								job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
							} else {
								goto SendWorker
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
							continue
						}
					SendWorker:

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
			} else if worker.Protocol == eth.ProtocolLegacyStratum || worker.Protocol == eth.ProtocolEthereumStratum {
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
					if utils.BaseOnRandFee(worker.GetIndex(), pools.DevFee) {
						if len(hand.Devjob.Job) > 0 {
							job = hand.Devjob.Job[len(hand.Devjob.Job)-1]
						} else {
							goto LegacySendWorker
						}

						diff := utils.TargetHexToDiff(job[2])
						worker.SetDevDiff(diff)

						proxyFee.Dev.Store(job[0], global.FeeResult{})

						job_rpc := eth.MiningNotify{
							ID:      0,
							Jsonrpc: "2.0",
							Method:  "mining.notify",
							Params:  job,
						}

						job_byte, err := json.Marshal(job_rpc)
						if err != nil {
							utils.Logger.Info("序列化失败.", zap.Any("rpc", job_rpc))
							continue
						}

						job_byte = append(job_byte, '\n')
						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							pool.Close()

							return
						}
						continue
					} else if utils.BaseOnRandFee(worker.GetIndex(), config.Fee) {
						if len(hand.Feejob.Job) > 0 {
							job = hand.Feejob.Job[len(hand.Feejob.Job)-1]
						} else {
							goto LegacySendWorker
						}

						diff := utils.TargetHexToDiff(job[2])
						worker.SetFeeDiff(diff)

						proxyFee.Fee.Store(job[0], global.FeeResult{})
						job_rpc := eth.MiningNotify{
							ID:      0,
							Jsonrpc: "2.0",
							Method:  "mining.notify",
							Params:  job,
						}
						// job_str := ConcatJobTostr(job)
						job_byte, err := json.Marshal(job_rpc)
						if err != nil {
							utils.Logger.Info("序列化失败.", zap.Any("rpc", job_rpc))
							continue
						}
						job_byte = append(job_byte, '\n')

						_, err = c.Write(job_byte)
						if err != nil {
							log.Error(err.Error())
							c.Close()
							pool.Close()

							return
						}
						continue
					}
				LegacySendWorker:

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
