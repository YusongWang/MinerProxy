package network

import (
	"crypto/tls"
	"log"
	"net"
)

func NewTls(crt string, key string) (ln net.Listener, err error) {
	//TODO check empty and give default cer.
	cer, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err = tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
