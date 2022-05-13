package handles

import (
	"io"
	"miner_proxy/global"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(io.ReadWriteCloser, *utils.Config, *global.Fee, string, *global.Worker) (io.ReadWriteCloser, error)
	OnMessage(io.ReadWriteCloser, *io.ReadWriteCloser, *utils.Config, *global.Fee, *[]byte, *global.Worker) ([]byte, error)
	OnClose(*global.Worker)
	SetLog(*zap.Logger)
}
