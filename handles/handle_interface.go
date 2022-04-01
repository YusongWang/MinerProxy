package handles

import (
	"net"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(string)
	OnMessage(net.Conn, []byte) ([]byte, error)
	OnClose()
	SetLog(*zap.Logger)
}
