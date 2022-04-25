package handles

import (
	"io"
	"miner_proxy/fee"
	"miner_proxy/pack"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(io.ReadWriteCloser, *utils.Config, *fee.Fee, string, *pack.Worker) (io.ReadWriteCloser, error)
	OnMessage(io.ReadWriteCloser, *io.ReadWriteCloser, *utils.Config, *fee.Fee, *[]byte, *pack.Worker) ([]byte, error)
	OnClose(*pack.Worker)
	SetLog(*zap.Logger)
}
