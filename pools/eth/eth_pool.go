package eth

import (
	"errors"
	"net"
	"strings"

	"miner_proxy/utils"
)

func NewPool(
	address string,
) (net.Conn, error) {
	if strings.HasPrefix(address, "tcp://") {
		address = strings.ReplaceAll(address, "tcp://", "")
		return utils.Tcp(address)
	} else if strings.HasPrefix(address, "ssl://") {
		address = strings.ReplaceAll(address, "ssl://", "")
		return utils.Tls(address)
	} else {
		return nil, errors.New("中转矿池: 不支持的协议类型: " + address)
	}
}
