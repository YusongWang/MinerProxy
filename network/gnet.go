package network

import (
	"fmt"
	"log"
	"miner_proxy/utils"

	"github.com/panjf2000/gnet/v2"
)

// TODO 框架原生不支持TLS。需要自己扩展。暂时接入了。
//
//

type NonBlockServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *NonBlockServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("NonBlock server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
	return gnet.None
}

func (es *NonBlockServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

func NewGnet(multicore bool, port int) {
	non_block := &NonBlockServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: multicore}
	utils.Logger.Error(gnet.Run(non_block, non_block.addr, gnet.WithMulticore(multicore)).Error())
}
