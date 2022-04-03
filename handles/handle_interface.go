package handles

import (
	"net"

	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(net.Conn, *utils.Config, string) (net.Conn, error)
	OnMessage(net.Conn, net.Conn, []byte) ([]byte, error)
	OnClose()
	SetLog(*zap.Logger)
}
