package serve

import (
	"bufio"
	"miner_proxy/fee"
	"miner_proxy/handles"
	"miner_proxy/utils"
	"net"

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
	return Serve{netln: netln, handle: handle, log: utils.Logger, config: config}
}

func (s *Serve) StartLoop() {
	for {
		// 循环接入所有客户端得到专线连接
		conn, err := s.netln.Accept()
		if err != nil {
			s.log.Error(err.Error())
			continue
		}

		s.log.Info("Tcp Accept Concent")
		s.handle.SetLog(s.log)

		var fee fee.Fee
		fee.Dev = make(map[string]bool)
		fee.Fee = make(map[string]bool)

		pool_net, err := s.handle.OnConnect(conn, s.config, &fee, conn.RemoteAddr().String())
		if err != nil {
			s.log.Warn(err.Error())
		}
		go s.serve(conn, pool_net, &fee)
	}
}

//接受请求
func (s *Serve) serve(conn net.Conn, pool net.Conn, fee *fee.Fee) {
	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			s.log.Error(err.Error())
			s.handle.OnClose()
			return
		}

		ret, err := s.handle.OnMessage(conn, pool, fee, buf)
		if err != nil {
			s.log.Error(err.Error())
			s.handle.OnClose()
			return
		}

		_, err = conn.Write(ret)
		if err != nil {
			s.log.Error(err.Error())
			s.handle.OnClose()
			return
		}
	}
}
