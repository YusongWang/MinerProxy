package main

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

func main() {
	defer ants.Release()
	runTimes := 100
	var wg sync.WaitGroup

	wg.Add(10)
	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	// p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
	// 	fmt.Println("form", i)
	// 	wg.Done()
	// })

	syncCalculateSum := func() {
		fmt.Println("Hello")
		wg.Done()
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}
	defer p.Release()
	wg.Wait()
}
