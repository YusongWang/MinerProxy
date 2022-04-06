package pack

import "sync"

type Job struct {
	Job  [][]string
	Lock sync.RWMutex
}
