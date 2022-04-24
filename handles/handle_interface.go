package handles

import (
	"io"
	"miner_proxy/fee"
	"miner_proxy/utils"

	"go.uber.org/zap"
)

type Handle interface {
	OnConnect(io.ReadWriteCloser, *utils.Config, *fee.Fee, string, *string) (io.ReadWriteCloser, error)
	OnMessage(io.ReadWriteCloser, *io.ReadWriteCloser, *utils.Config, *fee.Fee, *[]byte, *string) ([]byte, error)
	OnClose(*string)
	SetLog(*zap.Logger)
}
