package eth_stratum

import (
	"encoding/json"
	"errors"
	"sync"
)

type Job struct {
	Job  [][]string
	Lock sync.RWMutex
}

type JSONRpcReq struct {
	Id     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// Stratum
type JSONPushMessage struct {
	// FIXME: Temporarily add ID for Claymore compliance
	Id     interface{} `json:"id"`
	Result interface{} `json:"result"`
}

type JSONRpcResp struct {
	Id     interface{} `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}

type SubmitReply struct {
	Status string `json:"status"`
}

type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func EthStratumReq(data []byte) (JSONRpcReq, error) {
	var req JSONRpcReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

// Return Success
func EthSuccess(id json.RawMessage) (out []byte, err error) {
	rpc := &JSONRpcResp{
		Id:     id,
		Result: true,
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
