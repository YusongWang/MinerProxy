package serve

import (
	"bufio"
	"io"
	"miner_proxy/fee"
	"miner_proxy/handles"
	"miner_proxy/pack"
	pool "miner_proxy/pools"
	"miner_proxy/utils"
	"net"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

//TODO 定义接口传入以下参数
// handle
// global state
// pacakge
// rpc encode & decode
type Serve struct {
	config *utils.Config
	netln  net.Listener
	handle handles.Handle
	log    *zap.Logger
}

//主入口点。
// conn 支持TCP SSL
// handle 支持自定义handle.
//
func NewServe(
	netln net.Listener,
	handle handles.Handle,
	config *utils.Config,
) Serve {
	// 处理两个抽水矿工抽水率一致的问题
	// if utils.BaseFeeToIndex(config.Fee)%utils.BaseFeeToIndex(pool.DevFee) == 0 ||
	// 	utils.BaseFeeToIndex(pool.DevFee)%utils.BaseFeeToIndex(config.Fee) == 0 {
	// 	config.Fee += 0.1
	// }
	if utils.BaseFeeToIndex(config.Fee) == utils.BaseFeeToIndex(pool.DevFee) {
		config.Fee += 0.1
	}

	return Serve{netln: netln, handle: handle, log: utils.Logger, config: config}
}

func (s *Serve) StartLoop() {
	for {
		// 循环接入所有客户端得到专线连接
		conn, err := s.netln.Accept()
		if err != nil {
			s.log.Error(err.Error())
			return
		}

		s.handle.SetLog(s.log)

		var fee fee.Fee
		fee.Dev = make(map[string]bool)
		fee.Fee = make(map[string]bool)
		sessionId, err := uuid.NewRandom()
		if err != nil {
			s.log.Error(err.Error())
		}

		worker := pack.NewWorker("", "", sessionId.String(), conn.RemoteAddr().String())

		pool_net, err := s.handle.OnConnect(conn, s.config, &fee, conn.RemoteAddr().String(), worker)
		if err != nil {
			s.log.Warn(err.Error())
		}

		go s.serve(conn, &pool_net, &fee, worker)
	}
}

//接受请求
func (s *Serve) serve(conn io.ReadWriteCloser, pool *io.ReadWriteCloser, fee *fee.Fee, worker *pack.Worker) {
	defer func() {
		if x := recover(); x != nil {
			s.log.Info("Recover", zap.Any("err", x))
			return
		}
	}()

	reader := bufio.NewReader(conn)
	//TODO 处理通知所有线程结束任务

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {

			} else {
				s.log.Error(err.Error())
			}
			s.handle.OnClose(worker)
			return
		}

		ret, err := s.handle.OnMessage(conn, pool, s.config, fee, &buf, worker)
		if err != nil {
			if err == io.EOF {

			} else {
				s.log.Error(err.Error())
			}
			s.handle.OnClose(worker)
			return
		}

		// 兼容内部返回的情况
		if len(ret) > 0 {
			_, err = conn.Write(ret)
			if err != nil {
				if err == io.EOF {

				} else {
					s.log.Error(err.Error())
				}
				s.handle.OnClose(worker)
				return
			}
		}

	}
}
