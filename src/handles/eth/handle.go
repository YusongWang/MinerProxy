package eth

import (
	"encoding/json"
	"log"
	"miner_proxy/src/pack/eth"
	"net"
)

type Handle struct{}

func (hand *Handle) OnConnect(addr string) {
	log.Println("On Connect ", addr)
}

func (hand *Handle) OnMessage(c net.Conn, data []byte) (out []byte, err error) {
	log.Println(string(data))
	req, err := eth.EthStratumReq(data)
	if err != nil {
		log.Println(err)
		c.Close()
		return
	}

	switch req.Method {
	case "eth_submitLogin":
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			log.Println("Malformed stratum request params from", c.RemoteAddr().String())
			c.Close()
			return
		}

		// l := s.B.Listen()
		// go func() {
		// 	for {
		// 		select {
		// 		case job := <-l.Ch:
		// 			//fmt.Println(job)
		// 			c.Send(job)
		// 		}
		// 	}
		// }()

		// if req.Worker != "" {
		// 	s.Worker = req.Worker
		// } else {
		// 	p1 := strings.Split(params[0], ".")
		// }

		// reply, errReply := s.handleLoginRPC(cs, params, req.Worker)
		// if errReply != nil {
		// 	//return cs.sendTCPError(req.Id, errReply)
		// 	log.Println("Loign Error -1")
		// 	c.Close()
		// 	return
		// }
		//return cs.sendTCPResult(req.Id, reply)
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			log.Fatalln(err)
			c.Close()
			return
		}
		return
	case "eth_getWork":
		// reply, errReply := s.handleGetWorkRPC(cs)
		// if errReply != nil {
		// 	//return cs.sendTCPError(req.Id, errReply)
		// 	log.Println("Loign Error -1")
		// 	c.Close()
		// 	return
		// }
		// rpc := &eth.JSONRpcResp{
		// 	Id:      req.Id,
		// 	Version: "2.0",
		// 	Result:  true,
		// }

		// brpc, err := json.Marshal(rpc)
		// if err != nil {
		// 	log.Println(err)
		// 	c.Close()
		// 	return
		// }

		// log.Println("Ret", brpc)
		// out = append(brpc, '\n')
		return
	case "eth_submitWork":
		var params []string
		err = json.Unmarshal(req.Params, &params)
		if err != nil {
			log.Println("Malformed stratum request params from", c.RemoteAddr().String())
			return
		}
		//s.Remote <- params
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			log.Fatalln(err)
			c.Close()
			return
		}

		return
	case "eth_submitHashrate":
		// 直接返回
		out, err = eth.EthSuccess(req.Id)
		if err != nil {
			log.Fatalln(err)
			c.Close()
			return
		}

		return
	default:
		log.Println("KnownRpc")
		return
		// errReply := s.handleUnknownRPC(cs, req.Method)
		// return cs.sendTCPError(req.Id, errReply)
	}
	return
}

func (hand *Handle) OnClose() {

}
