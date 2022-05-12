package handles

import (
	"io"
	"miner_proxy/global"
	"miner_proxy/pack"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(io.ReadWriteCloser, *utils.Config, *global.Fee, string, *pack.Worker) (io.ReadWriteCloser, error)
	OnMessage(io.ReadWriteCloser, *io.ReadWriteCloser, *utils.Config, *global.Fee, *[]byte, *pack.Worker) ([]byte, error)
	OnClose(*pack.Worker)
	SetLog(*zap.Logger)
}
