package eth

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	pack "miner_proxy/src/pack/eth"
	"net"
	"runtime"
	"sync"

	broadcaster "github.com/tjgq/broadcast"
)

type EthStratumServer struct {
	Conn   net.Conn
	B      *broadcaster.Broadcaster
	Submit chan []string
}

func NewEthStratumServerSsl(
	address string,
	b *broadcaster.Broadcaster,
	submit chan []string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.B = b
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
	b *broadcaster.Broadcaster,
	submit chan []string,
) (EthStratumServer, error) {
	eth := EthStratumServer{}
	eth.B = b
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
	json_rpc := eth.ServerReq{
		ServerBaseReq: eth.ServerBaseReq{
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
func (eth *EthStratumServer) NotifyWorks(job) error {
	res, err := json.Marshal(job)
	if err != nil {
		return err
	}
	res = append(res, '\n')
	eth.B.Send(res)
	return nil
}

// 进行事件循环处理
func (eth *EthStratumServer) StartLoop() {
	var wg sync.WaitGroup

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

			var push eth.JSONPushMessage
			if err = json.Unmarshal([]byte(buf_str), &push); err == nil {
				eth.NotifyWorks(push)
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
