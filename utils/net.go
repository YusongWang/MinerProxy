package utils

import (
	"crypto/tls"
	"net"
)

func Tcp(address string) (net.Conn, error) {
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

func Tls(address string) (net.Conn, error) {
	cfg := tls.Config{}
	cfg.InsecureSkipVerify = true
	cfg.PreferServerCipherSuites = true
	conn, err := tls.Dial("tcp", address, &cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
