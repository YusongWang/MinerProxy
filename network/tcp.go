package network

import (
	"miner_proxy/utils"
	"net"
)

func NewTcp(addr string) (ln net.Listener, err error) {
	//TODO check empty and give default cer.
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	ln, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	return
}
