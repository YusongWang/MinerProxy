package serve

import (
	"bufio"
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
	netln  net.Listener
	handle handles.Handle
	log    *zap.Logger
}

//主入口点。
// conn 支持TCP SSL
// handle 支持自定义handle.
//
func NewServe(netln net.Listener, handle handles.Handle) Serve {
	return Serve{netln, handle, utils.Logger}
}

func (s *Serve) StartLoop() {
	for {
		// 循环接入所有客户端得到专线连接
		conn, err := s.netln.Accept()
		if err != nil {
			s.log.Error(err.Error())
			continue
		}
		s.log = s.log.With(zap.String("ip", conn.RemoteAddr().String()))
		s.log.Info("Tcp Accept Concent")
		s.handle.SetLog(s.log)
		s.handle.OnConnect(conn.RemoteAddr().String())
		go s.serve(conn)
	}
}

//接受请求
func (s *Serve) serve(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			s.log.Error(err.Error())
			s.handle.OnClose()
			return
		}

		ret, err := s.handle.OnMessage(conn, buf)
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