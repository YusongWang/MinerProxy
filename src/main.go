package main

import (
	"bufio"

	"log"
	"miner_proxy/src/cmd"
	"miner_proxy/src/handles/eth"

	"net"
)

func main() {
	//TODO init logger
	cmd.Execute()
	//TODO parse config
	// net, err := network.NewTcp(":38888")
	// if err != nil {
	// 	log.Panicln("can't bind to addr ", ":38888")
	// }
	// handle := eth.Handle{}
	// for {
	// 	//循环接入所有客户端得到专线连接
	// 	conn, err := net.Accept()
	// 	if err != nil {
	// 		log.Panicln(err.Error())
	// 		continue
	// 	}
	// 	log.Println("Tcp Accept Concent From ", conn.RemoteAddr())
	// 	handle.OnConnect(conn.RemoteAddr().String())
	// 	//开辟独立协程与该客聊天
	// 	go start_loop(conn, handle)
	// }
}

func start_loop(conn net.Conn, handle eth.Handle) {
	/// TODO 封装到serve层
	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}

		ret, err := handle.OnMessage(conn, buf)
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}

		_, err = conn.Write(ret)
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}
	}
}
