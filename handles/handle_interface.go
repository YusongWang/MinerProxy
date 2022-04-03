package handles

import (
	"net"

	"miner_proxy/fee"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(net.Conn, *utils.Config, *fee.Fee, string) (net.Conn, error)
	OnMessage(net.Conn, net.Conn, *fee.Fee, []byte) ([]byte, error)
	OnClose()
	SetLog(*zap.Logger)
}
