package eth

import (
	"encoding/json"
	"errors"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

// AuthorizeStat 认证状态
type AuthorizeStat uint8

const (
	StatConnected AuthorizeStat = iota
	StatSubScribed
	StatAuthorized
	StatDisconnected
	StatExit
)

// Stratum协议类型
type StratumProtocol uint8

const (
	// 未知协议
	ProtocolUnknown StratumProtocol = iota
	// ETHProxy 协议
	ProtocolETHProxy
	// NiceHash 的 EthereumStratum/1.0.0 协议
	ProtocolEthereumStratum
	// 传统 Stratum 协议
	ProtocolLegacyStratum
)

// NiceHash Ethereum Stratum Protocol 的协议类型前缀
const EthereumStratumPrefix = "ethereumstratum/"

// 响应中使用的 NiceHash Ethereum Stratum Protocol 的版本
const EthereumStratumVersion = "EthereumStratum/1.0.0"

//TODO 直接返回字符串。不用json解析

var ethsuccess = `{"id":`
var ethsuccess_end = `,"jsonrpc":"2.0","result":true}`

type JSONRpcReq struct {
	Id     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type JSONRpcReqType struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type StratumReq struct {
	JSONRpcReq
	Worker string `json:"worker"`
}

// Stratum
type JSONPushMessage struct {
	// FIXME: Temporarily add ID for Claymore compliance
	Id      int64       `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type JSONRpcResp struct {
	Id      json.RawMessage `json:"id"`
	Version string          `json:"jsonrpc"`
	Result  interface{}     `json:"result"`
	Error   interface{}     `json:"error,omitempty"`
}

type SubmitReply struct {
	Status string `json:"status"`
}

type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// REMOTE 远程服务器
type ServerBaseReq struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type ServerReq struct {
	ServerBaseReq
	Worker string `json:"worker"`
}

func EthStratumReq(data []byte) (StratumReq, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var req StratumReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

// Return Success
func EthSuccess(id int64) (out []byte, err error) {
	// rpc := &JSONRpcResp{
	// 	Id:      id,
	// 	Version: "2.0",
	// 	Result:  true,
	// }
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// out, err = json.Marshal(rpc)
	// if err != nil {
	// 	return nil, err
	// }

	out = []byte(ethsuccess + strconv.Itoa(int(id)) + ethsuccess_end)
	out = append(out, '\n')
	return
}

func EthError(id json.RawMessage, code int32, msg string) ([]byte, error) {
	return nil, errors.New("TODO")
}

type JSONRPCArray []interface{}

///{"id":1,"method":"mining.subscribe","params":["MinerName/1.0.0",""]}
type MiningNotify struct {
	ID      int          `json:"id"`
	Jsonrpc string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  JSONRPCArray `json:"params"`
}
