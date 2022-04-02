package handles

import (
	"net"

	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(net.Conn, *utils.Config, string) error
	OnMessage(net.Conn, []byte) ([]byte, error)
	OnClose()
	SetLog(*zap.Logger)
}
