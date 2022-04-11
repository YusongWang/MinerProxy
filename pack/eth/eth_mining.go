package eth

import (
	"encoding/json"
	"errors"
)

type JSONRpcReq struct {
	Id     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
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

	var req StratumReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

// Return Success
func EthSuccess(id json.RawMessage) (out []byte, err error) {
	rpc := &JSONRpcResp{
		Id:      id,
		Version: "2.0",
		Result:  true,
	}

	out, err = json.Marshal(rpc)
	if err != nil {
		return nil, err
	}
	return
}

func EthError(id json.RawMessage, code int32, msg string) ([]byte, error) {
	return nil, errors.New("TODO")
}
