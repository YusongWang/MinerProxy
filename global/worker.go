package global

import (
	"miner_proxy/pack"
	"sync"
)

type OnlineWorkers struct {
	Workers map[string]*pack.Worker
	sync.Mutex
}

var GonlineWorkers = new(OnlineWorkers)

func init() {
	GonlineWorkers.Lock()
	GonlineWorkers.Workers = make(map[string]*pack.Worker)
	GonlineWorkers.Unlock()
}
