package global

import (
	"miner_proxy/pack"
	"sync"
)

type OnlineWorkers struct {
	Workers map[string]*pack.Worker
	sync.Mutex
}

var GonlineWorkers OnlineWorkers
