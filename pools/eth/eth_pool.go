package eth

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	pack "miner_proxy/pack/eth"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type EthStratumServer struct {
	Conn   net.Conn
	Job    *pack.Job
	Submit chan []string
}

func NewEthStratumServerSsl(
	address string,
	job *pack.Job,
	submit chan []string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.Job = job
	eth.Submit = submit

	cfg := tls.Config{}
	cfg.InsecureSkipVerify = true
	cfg.PreferServerCipherSuites = true
	conn, err := tls.Dial("tcp", address, &cfg)
	//conn, err := net.Dial("tcp", address)
	if err != nil {
		return eth, err
	}
	eth.Conn = conn
	return eth, nil
}

func NewEthStratumServerTcp(
	address string,
	job *pack.Job,
	submit chan []string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.Job = job
	eth.Submit = submit
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return eth, err
	}
	eth.Conn = conn
	return eth, nil
}

// 用自定义钱包进行登陆
func (eth *EthStratumServer) Login(wallet string) error {
	//eth_mining.EthStratumReq()
	var a []string
	a = append(a, wallet)
	a = append(a, "x")
	login := pack.ServerReq{
		ServerBaseReq: pack.ServerBaseReq{
			Id:     0,
			Method: "eth_submitLogin",
			Params: a,
		},
		Worker: "P1",
	}

	res, err := json.Marshal(login)
	if err != nil {
		return err
	}

	log.Println(string(res))

	write := append(res, '\n')
	len, err := eth.Conn.Write(write)
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

// 提交工作量证明
func (eth *EthStratumServer) SubmitJob(job []string) error {
	json_rpc := pack.ServerReq{
		ServerBaseReq: pack.ServerBaseReq{
			Id:     40,
			Method: "eth_submitWork",
			Params: job,
		},
		Worker: "Default",
	}

	log.Println("给服务器提交工作量证明", json_rpc)
	res, err := json.Marshal(json_rpc)
	if err != nil {
		log.Println("Json Marshal Error ", err)
		return err
	}

	ret := append(res, '\n')
	_, err = eth.Conn.Write(ret)
	if err != nil {
		return err
	}

	return nil
}

// bradcase 当前工作
func (eth *EthStratumServer) NotifyWorks(job []string) error {
	eth.Job.Lock.Lock()
	eth.Job.Job = append(eth.Job.Job, job)
	eth.Job.Lock.Unlock()
	return nil
}

// 进行事件循环处理
func (eth *EthStratumServer) StartLoop() {
	var wg sync.WaitGroup
	// go func() {
	// 	for {
	// 		eth.Job.Lock.RLock()
	// 		log.Println(eth.Job.Job)
	// 		eth.Job.Lock.RUnlock()
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			buf_str, err := bufio.NewReader(eth.Conn).ReadString('\n')
			if err != nil {
				log.Println("远程已经关闭")
				log.Println(err)
				eth.Conn.Close()
				return
			}
			var push pack.JSONPushMessage
			if err = json.Unmarshal([]byte(buf_str), &push); err == nil {
				if result, ok := push.Result.(bool); ok {
					//增加份额
					if result == true {
						// TODO
						utils.Logger.Info("有效份额", zap.Any("RPC", buf_str))
					} else {
						utils.Logger.Warn("无效份额", zap.Any("RPC", buf_str))
					}
				} else if list, ok := push.Result.([]interface{}); ok {
					job := make([]string, len(list))
					for i, arg := range list {
						job[i] = arg.(string)
					}
					eth.NotifyWorks(job)
				} else {
					//TODO
				}
			} else {
				log.Println(err)
			}

			log.Println(buf_str)
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for {
				select {
				case job := <-eth.Submit:
					err := eth.SubmitJob(job)
					if err != nil {
						log.Fatalln("提交工作量证明失败")
					}
					log.Println(job)
				}
			}

		}()
	}

	wg.Wait()
}
