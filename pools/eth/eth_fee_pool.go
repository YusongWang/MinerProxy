package eth

import (
	"bufio"
	"errors"
	"io"
	"log"
	"miner_proxy/global"
	ethpack "miner_proxy/pack/eth"
	"miner_proxy/utils"
	"os"
	"strings"
	"sync"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type setNoDelayer interface {
	SetNoDelay(bool) error
}

type EthStratumServer struct {
	Conn     *io.ReadWriteCloser
	Job      *[]global.Job
	Submit   chan []byte
	PoolAddr string
	Wallet   string
	Worker   string
}

func New(
	address string,
	job *[]global.Job,
	submit chan []byte,
) (EthStratumServer, error) {
	if strings.HasPrefix(address, "tcp://") {
		address = strings.ReplaceAll(address, "tcp://", "")
		return newEthStratumServerTcp(address, job, submit, address)
	} else if strings.HasPrefix(address, "ssl://") {
		address = strings.ReplaceAll(address, "ssl://", "")
		return newEthStratumServerSsl(address, job, submit, address)
	} else {
		return EthStratumServer{}, errors.New("不支持的协议类型: " + address)
	}
}

func newEthStratumServerSsl(
	address string,
	job *[]global.Job,
	submit chan []byte,
	pool string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.Job = job
	eth.Submit = submit
	eth.PoolAddr = "ssl://" + address

	conn, err := utils.Tls(address)
	if err != nil {
		return eth, err
	}
	eth.Conn = &conn
	// cfg := tls.Config{}
	// cfg.InsecureSkipVerify = true
	// cfg.PreferServerCipherSuites = true
	// var err error
	// *eth.Conn, err = tls.Dial("tcp", address, &cfg)
	// if err != nil {
	// 	return eth, err
	// }
	if c, ok := (*eth.Conn).(setNoDelayer); ok {
		c.SetNoDelay(true)
	}
	return eth, nil
}

func newEthStratumServerTcp(
	address string,
	job *[]global.Job,
	submit chan []byte,
	pool string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.Job = job
	eth.Submit = submit
	eth.PoolAddr = "tcp://" + address
	// tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	// if err != nil {
	// 	return eth, err
	// }

	conn, err := utils.Tcp(address)
	eth.Conn = &conn
	if err != nil {
		return eth, err
	}

	if c, ok := (*eth.Conn).(setNoDelayer); ok {
		c.SetNoDelay(true)
	}

	return eth, nil
}

// 用自定义钱包进行登陆
func (eth *EthStratumServer) Login(wallet string, worker string) error {
	//eth_mining.EthStratumReq()
	if worker == "" {
		return errors.New("矿工名称不能为空！")
	}
	eth.Worker = worker
	eth.Wallet = wallet

	var a []string
	a = append(a, wallet)
	a = append(a, "x")
	login := ethpack.ServerReq{
		ServerBaseReq: ethpack.ServerBaseReq{
			Id:     0,
			Method: "eth_submitLogin",
			Params: a,
		},
		Worker: worker,
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	res, err := json.Marshal(login)
	if err != nil {
		return err
	}

	write := append(res, '\n')
	len, err := (*eth.Conn).Write(write)
	if err != nil {
		log.Println("Socket Close", err)
		return err
	}

	if len <= 0 {
		log.Println("Socket Close len :", len)
		return errors.New("Socket Close len")
	}

	return nil
}

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

var package_head = `{"id":40,"method":"eth_submitWork","params":`
var package_middle = `,"worker":"`
var package_end = `"}`

// 提交工作量证明
func (eth *EthStratumServer) SubmitJob(job []byte) error {
	//str := ConcatJobTostr(job)
	var builder strings.Builder
	builder.WriteString(package_head)
	builder.WriteString(string(job))
	builder.WriteString(package_middle)
	builder.WriteString(eth.Worker)
	builder.WriteString(package_end)
	builder.WriteByte('\n')
	json_rpc := builder.String()

	//utils.Logger.Info("给服务器提交工作量证明", zap.Any("RPC", json_rpc))

	_, err := (*eth.Conn).Write([]byte(json_rpc))
	if err != nil {
		return err
	}

	return nil
}

// bradcase 当前工作
func (eth *EthStratumServer) NotifyWorks(buf *[]byte) error {

	job_len := len(*eth.Job)
	if job_len >= 1000 {
		*eth.Job = (*eth.Job)[500:job_len]
	}

	job_diff, err := jsonparser.GetString(*buf, "result", "[2]")
	if err != nil {
		utils.Logger.Error("无法解析result Diff字段")
	}

	job_id, err := jsonparser.GetString(*buf, "result", "[0]")
	if err != nil {
		utils.Logger.Error("无法解析result JobId字段")
	}

	target, err := jsonparser.GetString(*buf, "result", "[1]")
	if err != nil {
		utils.Logger.Error("无法解析result Target字段")
	}

	job_byte := append(*buf, '\n')
	j := global.Job{
		Target: target,
		JobId:  job_id,
		Diff:   job_diff,
		Job:    job_byte,
	}

	*eth.Job = append(*eth.Job, j)
	return nil
}

// 进行事件循环处理
func (eth *EthStratumServer) StartLoop() {
	var wg sync.WaitGroup

	log := utils.Logger.With(zap.String("Worker", eth.Worker))
	wg.Add(1)
	go func() {
		reader := bufio.NewReader(*eth.Conn)
		var buf []byte
		var err error
		defer wg.Done()
		for {
			buf, err = reader.ReadBytes('\n')
			if err != nil {
				log.Info("矿池关闭->  尝试重新连接")
				log.Error(err.Error())
				temp, err := New(eth.PoolAddr, eth.Job, eth.Submit)
				if err != nil {
					log.Info("矿池关闭->  尝试重新连接失败 ")
					log.Error(err.Error())
					os.Exit(1)
				}
				err = temp.Login(eth.Wallet, eth.Worker)
				if err != nil {
					log.Info("矿池关闭->  尝试重新连接失败 ")
					log.Error(err.Error())
					os.Exit(1)
				}

				*eth.Conn = *temp.Conn
				reader = bufio.NewReader(*eth.Conn)
				continue
			}

			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			var push ethpack.JSONPushMessage
			if err = json.Unmarshal(buf, &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result {
						// TODO
						//log.Info("有效份额", zap.Any("RPC", buf_str))
					} else {
						log.Warn("无效份额", zap.String("RPC", string(buf)))
					}
				} else if _, ok := push.Result.([]interface{}); ok {
					// job := make([]string, len(list))
					// for i, arg := range list {
					// 	job[i] = arg.(string)
					// }
					eth.NotifyWorks(&buf)
				} else {
					//TODO
				}
			} else {
				log.Error(err.Error())
			}
		}
	}()

	//TODO 调试这里的最优化接受携程数量
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		for {
	// 			select {
	// 			case job := <-eth.Submit:
	// 				err := eth.SubmitJob(job)
	// 				if err != nil {
	// 					log.Warn("提交工作量证明失败")
	// 				}
	// 			}
	// 		}
	// 	}()
	// }

	wg.Wait()
}
