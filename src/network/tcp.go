package network

import (
	"log"
	"net"
)

func NewTcp(addr string) (ln net.Listener, err error) {
	//TODO check empty and give default cer.
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Println(err.Error())
	}
	ln, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
