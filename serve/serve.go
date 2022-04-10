package serve

import (
	"bufio"
	"miner_proxy/fee"
	"miner_proxy/handles"
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
	if int(config.Fee*10)%int(pool.DevFee*10) == 0 || int(pool.DevFee*10)%int(config.Fee*10) == 0 {
		config.Fee += config.Fee + 0.1
	}

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
		bid, err := uuid.NewUUID()
		if err != nil {
			s.log.Error(err.Error())
		}
		id := bid.String()
		pool_net, err := s.handle.OnConnect(conn, s.config, &fee, conn.RemoteAddr().String(), &id)
		if err != nil {
			s.log.Warn(err.Error())
		}
		go s.serve(conn, pool_net, &fee, &id)
	}
}

//接受请求
func (s *Serve) serve(conn net.Conn, pool net.Conn, fee *fee.Fee, id *string) {

	reader := bufio.NewReader(conn)
	// dev_bufio := s.Handle.GiveDevbufio()
	// fee_bufio := s.Handle.GiveFeebufio()
	//TODO
	//w3 := bufio.NewWriter(conn);

	// go func(ww * bufio.Writer) {
	// 	for ;;  {
	// 		ww.WriteString("hi33333333333333333333333333333333333333333333\r\n");
	// 		ww.Flush();
	// 	}
	// }(w3)

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			s.log.Error(err.Error())

			s.handle.OnClose(id)
			return
		}

		go func(buf []byte) {
			ret, err := s.handle.OnMessage(conn, pool, fee, buf, id)
			if err != nil {
				s.log.Error(err.Error())
				s.handle.OnClose(id)
				return
			}

			// 兼容内部返回的情况
			if len(ret) > 0 {
				_, err = conn.Write(ret)
				if err != nil {
					s.log.Error(err.Error())
					s.handle.OnClose(id)
					return
				}
			}
		}(buf)
	}
}
