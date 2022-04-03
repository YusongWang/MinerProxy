package eth

import (
	"net"

	"go.uber.org/zap"
)

type HandlePool struct {
	log *zap.Logger
}

func (hand *HandlePool) OnConnect(addr string) {
	hand.log.Info("On Connect")
}

func (hand *HandlePool) OnMessage(c net.Conn, data []byte) (out []byte, err error) {

	return
}

func (hand *HandlePool) OnClose() {
	hand.log.Info("OnClose !!!!!")
}

func (hand *HandlePool) SetLog(log *zap.Logger) {
	hand.log = log
}
