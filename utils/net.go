package utils

import (
	"crypto/tls"
	"io"
	"net"
)

type setNoDelayer interface {
	SetNoDelay(bool) error
}

func Tcp(address string) (io.ReadWriteCloser, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	conn.SetNoDelay(true)

	return conn, nil
}

func Tls(address string) (io.ReadWriteCloser, error) {
	var rw io.ReadWriteCloser

	cfg := tls.Config{}
	cfg.InsecureSkipVerify = true
	cfg.PreferServerCipherSuites = true
	rw, err := tls.Dial("tcp", address, &cfg)
	if err != nil {
		return nil, err
	}
	if c, ok := rw.(setNoDelayer); ok {
		c.SetNoDelay(true)
	}

	return rw, nil
}
