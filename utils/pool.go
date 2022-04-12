package utils

import (
	"errors"
	"io"
	"strings"
)

func NewPool(
	address string,
) (io.ReadWriteCloser, error) {
	if strings.HasPrefix(address, "tcp://") {
		address = strings.ReplaceAll(address, "tcp://", "")
		return Tcp(address)
	} else if strings.HasPrefix(address, "ssl://") {
		address = strings.ReplaceAll(address, "ssl://", "")
		return Tls(address)
	} else {
		return nil, errors.New("中转矿池: 不支持的协议类型: " + address)
	}
}
