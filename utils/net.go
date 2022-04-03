package utils

import (
	"crypto/tls"
	"net"
)

func Tcp(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
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
